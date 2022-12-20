package components

type JobResp struct {
	// Component instance ID.
	ComponentId string `json:"component_id"`
	// Job ID.
	JobId string `json:"job_id"`
}
type CreateOpts struct {
	// Enterprise project id
	EnterpriseProjectId string `json:"enterprise_project_id,omitempty"`
	// Application component name.
	// The value can contain 2 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	// It must start with a letter and end with a letter or digit.
	Name string `json:"name" required:"true"`
	// The component version
	Version string `json:"version",omitempty`
	// Environment id
	EnvironmentId string `json:"environment_id" required:"true"`
	// workload name
	WorkloadName string `json:"workload_name,omitempty"`
	// Application id
	ApplicationId string `json:"application_id" required:"true"`
	// Description.
	// The value can contain up to 128 characters.
	Description string `json:"description,omitempty"`
	// JVM config.
	JVMOpts string `json:"jvm_opts,omitempty"`
	// Technology stack
	RuntimeStack *RuntimeStack `json:"runtime_stack" required:"true"`
	// Tomcat config
	TomcatOpts *TomcatOptions `json:"tomcat_opts,omitempty"`
	// Component builder.
	Builder *Builder `json:"build,omitempty"`
	// Source of the code or software package.
	Source *Source `json:"source,omitempty"`
	// Workload kind
	WorkloadKind string `json:"workload_kind,omitempty"`

	LimitCPU int `json:"limit_cpu,omitempty"`

	LimitMemory int `json:"limit_memory,omitempty"`

	RequestCPU int `json:"request_cpu,omitempty"`

	RequestMemory int `json:"request_memory,omitempty"`
	// default 1
	Replica int `json:"replica,omitempty"`

	ServiceName string `json:"service_name,omitempty"`

	Ports []Port `json:"ports,omitempty"`
	// External network access.
	ExternalAccesses []ExternalAccess `json:"external_accesses,omitempty"`

	Acceleration *Acceleration `json:"acceleration,omitempty"`
	// Environment variable.
	EnvVariables []EnvVariable `json:"env,omitempty"`

	Commands *Command `json:"command" required:"true"`
	// Data storage configuration.
	Storages []Storage `json:"storage,omitempty"`
	// Post-start processing.
	PostStart *Process `json:"post_start,omitempty"`
	// Pre-stop processing.
	PreStop *Process `json:"pre_stop,omitempty"`
	// Affinity.
	Affinity []Affinity `json:"affinity,omitempty"`
	// Anti-affinity.
	AntiAffinity []Affinity `json:"anti_affinity,omitempty"`

	SecurityContext *SecurityContext `json:"security_context,omitempty"`

	DNSPolicy string `json:"dns_policy,omitempty"`

	DNSConfig *DNSConfig `json:"dns_config,omitempty"`
	// Policy list of log collection.
	LogCollectionPolicies []LogCollectionPolicy `json:"logs,omitempty"`
	// Component liveness probe.
	LivenessProbe *ProbeDetail `json:"livenessProbe,omitempty"`
	// Component service probe.
	ReadinessProbe *ProbeDetail `json:"readinessProbe,omitempty"`

	CustomMetric *CustomMetric `json:"custom_metric,omitempty"`
	// Deployed resources.
	ReferResources []ReferResource `json:"refer_resources,omitempty"`

	Labels []string `json:"labels,omitempty"`

	UpdateStrategy *UpdateStrategy `json:"update_strategy,omitempty"`

	DeployStrategy *DeployStrategy `json:"deploy_strategy,omitempty"`
}
