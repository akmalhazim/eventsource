package eventsource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/akmalhazim/eventsource"
	"github.com/akmalhazim/eventsource/aggregatestore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	UserAggregate = eventsource.AggregateType("User")
)

var (
	Active     = UserStatus("Active")
	Disabled   = UserStatus("Disabled")
	Unverified = UserStatus("Unverified")
)

type UserStatus string

type User struct {
	*eventsource.BaseAggregate
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
	Status   UserStatus
}

func (user *User) AggregateID() uuid.UUID {
	return user.ID
}

func (user *User) AggregateType() eventsource.AggregateType {
	return UserAggregate
}

func (user *User) HandleEvent(event eventsource.Event) {
	switch evt := event.(type) {
	case *UserRegistered:
		user.Name = evt.Name
		user.Email = evt.Email
		user.Password = evt.Password
	case *UserVerified:
		user.Status = Active
	}
}

type UserRegistered struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (UserRegistered) EventType() eventsource.EventType {
	return eventsource.EventType("UserRegistered")
}

type UserVerified struct{}

func (UserVerified) EventType() eventsource.EventType {
	return eventsource.EventType("UserVerified")
}

func init() {
	eventsource.RegisterAggregate(func(id uuid.UUID) eventsource.Aggregate {
		return &User{
			BaseAggregate: &eventsource.BaseAggregate{},
			ID:            id,
			Status:        Unverified,
		}
	})

	eventsource.RegisterEvent(func() eventsource.Event {
		return &UserRegistered{}
	})

	eventsource.RegisterEvent(func() eventsource.Event {
		return &UserVerified{}
	})
}

func TestCreateUserAggregate(t *testing.T) {
	t.SkipNow()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	db := client.Database("eventsource_test")

	store := aggregatestore.NewMongoAggregateStore(db)
	aggregate, err := store.Load(context.TODO(), UserAggregate, uuid.MustParse("c5511d2b-f200-484b-9ff2-9abd1940c462"))
	if err != nil {
		fmt.Println(err)
	}

	assert.Nil(t, err, "error is not nil")
	assert.NotNil(t, aggregate, "user is nil")

	user := aggregate.(*User)
	user.AppendEvent(&UserRegistered{
		Name:     "Amira Syahirah",
		Email:    "nuramirasyahirah@gmail.com",
		Password: "password",
	})

	err = store.Save(context.TODO(), user)
	assert.Nil(t, err, "error is not nil")
	assert.NotNil(t, user, "user is nil")
}

func TestVerifyUserAggregate(t *testing.T) {
	t.SkipNow()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	db := client.Database("eventsource_test")
	store := aggregatestore.NewMongoAggregateStore(db)
	aggregate, err := store.Load(context.TODO(), UserAggregate, uuid.MustParse("c5511d2b-f200-484b-9ff2-9abd1940c462"))
	assert.Nil(t, err)
	assert.NotNil(t, aggregate)

	user := aggregate.(*User)
	user.AppendEvent(&UserVerified{})

	err = store.Save(context.TODO(), user)
	assert.Nil(t, err)
	assert.Equal(t, Active, user.Status)
}

func TestLoadUserAggregate(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	db := client.Database("eventsource_test")

	store := aggregatestore.NewMongoAggregateStore(db)
	aggregate, err := store.Load(context.TODO(), UserAggregate, uuid.MustParse("c5511d2b-f200-484b-9ff2-9abd1940c462"))
	assert.Nil(t, err, "aggregate can't be nil")

	user := aggregate.(*User)
	assert.NotNil(t, user)
	assert.Equal(t, "Amira Syahirah", user.Name, "user name is not correct")
	assert.Equal(t, "nuramirasyahirah2@gmail.com", user.Email, "user email is not correct")
	assert.Equal(t, Active, user.Status, "user is not active")

	fmt.Println(user)
}
