package examples

type Spec struct {
	CronSpec string
	Image    string
	Replicas int
}
type CronTab struct {
	Spec Spec
}
