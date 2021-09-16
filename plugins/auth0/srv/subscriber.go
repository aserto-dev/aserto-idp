package srv

import (
	"encoding/json"
	"log"
	"math"
	"sync"
	"time"

	"github.com/aserto-dev/aserto-idp/plugins/auth0/config"
	api "github.com/aserto-dev/go-grpc/aserto/api/v1"
	"gopkg.in/auth0.v5"
	"gopkg.in/auth0.v5/management"
)

const (
	maxBatchSize = int64(500 * 1024)
)

// Subscriber
type Subscriber struct {
	cfg   *config.Auth0Config
	count int
}

// NewSubscriber returns Auth0 producer instance.
func NewSubscriber(cfg *config.Auth0Config) *Subscriber {
	return &Subscriber{
		cfg: cfg,
	}
}

func (p *Subscriber) Count() int {
	return p.count
}

// Subscriber func.
func (p *Subscriber) Subscriber(s <-chan *api.User, errc chan<- error, done chan<- bool) {
	mgnt, err := management.New(
		p.cfg.Domain,
		management.WithClientCredentials(
			p.cfg.ClientID,
			p.cfg.ClientSecret,
		))
	if err != nil {
		errc <- err
		return
	}

	c, err := mgnt.Connection.ReadByName("Username-Password-Authentication")
	if err != nil {
		errc <- err
	}

	connectionID := auth0.StringValue(c.ID)

	var users []map[string]interface{}
	var jobs []management.Job
	var totalSize int64
	var wg sync.WaitGroup

	for u := range s {
		user, err := TransformToAuth0(u)
		if err != nil {
			errc <- err
		}
		userMap, size, err := structToMap(user)
		if err != nil {
			errc <- err
		}

		if totalSize+size < maxBatchSize {
			users = append(users, userMap)
		} else {
			// create job with users so far
			job := &management.Job{
				ConnectionID:        auth0.String(connectionID),
				Upsert:              auth0.Bool(true),
				SendCompletionEmail: auth0.Bool(false),
				Users:               users,
			}
			jobs = append(jobs, *job)
			users = make([]map[string]interface{}, 0)
			users = append(users, userMap)
		}
	}

	// create job for remaining users
	job := &management.Job{
		ConnectionID:        auth0.String(connectionID),
		Upsert:              auth0.Bool(true),
		SendCompletionEmail: auth0.Bool(false),
		Users:               users,
	}
	jobs = append(jobs, *job)

	for _, job := range jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = mgnt.Job.ImportUsers(&job)
			if err != nil {
				errc <- err
			}
			jobID := auth0.StringValue(job.ID)
			log.Printf("waiting for import job %s to finish ...", jobID)
			for {
				j, err := mgnt.Job.Read(jobID)
				if err != nil {
					errc <- err
				}
				if *j.Status != "pending" {
					break
				}
				time.Sleep(1 * time.Second)
			}
			log.Printf("Finished %s", jobID)
		}()
	}
	wg.Wait()
	done <- true
}

func structToMap(in interface{}) (map[string]interface{}, int64, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, 0, err
	}
	res := make(map[string]interface{})
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, 0, err
	}
	size := int64(len(data))
	return res, size, nil
}

func eventSize(u api.User) int64 {
	if b, err := json.Marshal(u); err == nil {
		return int64(len(b))
	}
	return math.MaxInt64
}
