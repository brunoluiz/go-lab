package healthz

import "github.com/brunoluiz/go-lab/lib/httpx"

// add some checks on instance creation
// h, _ := health.New(health.WithComponent(health.Component{
// 	Name:    "myservice",
// 	Version: "v1.0",
// }), health.WithChecks(health.Config{
// 	Name:      "rabbitmq",
// 	Timeout:   time.Second * 5,
// 	SkipOnErr: true,
// 	Check: func(ctx context.Context) error {
// 		// rabbitmq health check implementation goes here
// 		return nil
// 	}}, health.Config{
// 	Name: "mongodb",
// 	Check: func(ctx context.Context) error {
// 		// mongo_db health check implementation goes here
// 		return nil
// 	},
// },
// ))
//
// // and then add some more if needed
// h.Register(health.Config{
// 	Name:      "mysql",
// 	Timeout:   time.Second * 2,
// 	SkipOnErr: false,
// 	Check: healthMysql.New(healthMysql.Config{
// 		DSN: "test:test@tcp(0.0.0.0:31726)/test?charset=utf8",
// 	}),
// })
//
// http.Handle("/status", h.Handler())
// http.ListenAndServe(":3000", nil)

func New() (*httpx.Server, *health.Health) {
}
