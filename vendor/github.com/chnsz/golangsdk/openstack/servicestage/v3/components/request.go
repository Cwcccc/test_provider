package components

import (
"github.com/chnsz/golangsdk"
"github.com/chnsz/golangsdk/pagination"
)

var requestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpts is the structure required by the Create method to create a new component.
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

	Acceleration Acceleration `json:"acceleration,omitempty"`
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

	SecurityContext SecurityContext `json:"security_context,omitempty"`

	DNSPolicy string `json:"dns_policy,omitempty"`

	DNSConfig DNSConfig `json:"dns_config,omitempty"`
	// Policy list of log collection.
	LogCollectionPolicies []LogCollectionPolicy `json:"logs,omitempty"`
	// Component liveness probe.
	LivenessProbe *ProbeDetail `json:"livenessProbe,omitempty"`
	// Component service probe.
	ReadinessProbe *ProbeDetail `json:"readinessProbe,omitempty"`
	
	CustomMetric CustomMetric `json:"custom_metric,omitempty"`
	// Deployed resources.
	ReferResources []ReferResource `json:"refer_resources,omitempty"`

	Labels []string `json:"labels,omitempty"`

	UpdateStrategy UpdateStrategy `json:"update_strategy,omitempty"`

	DeployStrategy DeployStrategy `json:"deploy_strategy,omitempty"`
}

type RuntimeStack struct {
	// Technology stack name
	Name string `json:"name" required:"true"`
	// Technology stack version
	Version string `json:"version,omitempty"`
	// Deploy mode. Value: container, virturalMachicne
	DeployMode string `json:"deploy_mode" required:"true"`
	// Type
	Type string `json:"type" required:"true"`
}

type TomcatOptions struct {
	ServerXML string `json:"server_xml,omitempty"`
}

// Builder is the component builder, the configuration details refer to 'Parameter'.
type Builder struct {
	// Compilation command. By default:
	// When build.sh exists in the root directory, the command is ./build.sh.
	// When build.sh does not exist in the root directory, the command varies depending on the operating system (OS). Example:
	// Java and Tomcat: mvn clean package
	// Nodejs: npm build
	BuildCmd string `json:"build_cmd,omitempty"`
	// Address of the Docker file. By default, the Docker file is in the root directory (./).
	DockerfilePath string `json:"dockerfile_path,omitempty"`
	// Build archive organization. Default value: cas_{project_id}.
	ArtifactNamespace string `json:"artifact_namespace,omitempty"`
	// key indicates the key of the tag, and value indicates the value of the tag.
	UsePublicCluster bool `json:"use_public_cluster,omitempty"`
	// key indicates the key of the tag, and value indicates the value of the tag.
	NodeLabelSelector map[string]interface{} `json:"node_label_selector,omitempty"`
	//
	EnvironmentId string `json:"environment_id,omitempty"`
}

// Source is an object to specified the source information of Open-Scoure codes or package storage.
type Source struct {
	// The parameters of artifact are as follows:
	// Type. Option: source code or artifact software package.
	Kind string `json:"kind" required:"true"`
	// Address of the software package or source code.
	Url string `json:"url,omitempty"`
	//
	Version string `json:"version,omitempty"`
	// Storage mode. Value: swr or obs.
	Storage string `json:"storage,omitempty"`

	// The parameters of code are as follows:
	// Code repository. Value: GitHub, GitLab, Gitee, or Bitbucket.
	RepoType string `json:"repo_type,omitempty"`
	// Code repository URL. Example: https://github.com/example/demo.git.
	RepoUrl string `json:"repo_url,omitempty"`
	// Authorization name, which can be obtained from the authorization list.
	RepoAuth string `json:"repo_auth,omitempty"`
	// The code's organization. Value: GitHub, GitLab, Gitee, or Bitbucket.
	RepoNamespace string `json:"repo_namespace,omitempty"`
	// Code branch or tag. Default value: master.
	RepoRef string `json:"repo_ref,omitempty"`
}

type Port struct {
	Name string `json:"name,omitempty"`
	Port int `json:"port,omitempty"`
	TargetPort int `json:"target_port,omitempty"`
}

// ExternalAccess is an object that specifies the configuration of the external IP access.
type ExternalAccess struct {
	// Protocol. Value: http or https.
	Protocol string `json:"protocol,omitempty"`
	// Access address.
	Address string `json:"address,omitempty"`
	// Port number.
	ForwardPort int `json:"forward_port,omitempty"`
	ID string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type Acceleration struct {
	ClaimName string `json:"claim_name,Acceleration"`
}

type Command struct {
	Command string `json:"command,omitempty"`
	Args []string `json:"args,omitempty"`
}

// Storage is an object that specifies the data storage.
type Storage struct {
	// Storage type. Value:
	// HostPath: host path mounting.
	// EmptyDir: temporary directory mounting.
	// ConfigMap: configuration item mounting.
	// Secret: secret volume mounting.
	// PersistentVolumeClaim: cloud storage mounting.
	Type string `json:"type" required:"true"`
	// Storage parameter.
	Parameters *StorageParams `json:"parameters" required:"true"`
	// Directory mounted to the container.
	Mounts *Mount `json:"mounts" required:"true"`
}

type StorageParams struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
	DefaultMode string `json:"default_mode,omitempty"`
	Medium string `json:"medium,omitempty"`
}

type Mount struct {
	Path string `json:"path,omitempty"`
	SubPath string `json:"sub_path,omitempty"`
	Readonly bool `json:"readonly,omitempty"`
}

// Configuration is an object that specifies the environment variable for the component instance.
type EnvVariable struct {
	// Environment variable name.
	// The value contains 1 to 64 characters, including letters, digits, underscores (_), hyphens (-), and dots (.),
	// and cannot start with a digit.
	Name string `json:"name" required:"true"`
	// Environment variable value.
	Value string `json:"value" required:"true"`
	Inner bool `json:"inner" required:"true"`
	ValueFrom ValueFrom `json:"value_from,omitempty"`
}

type ValueFrom struct {
	ReferenceType string `json:"reference_type,omitempty"`
	Name string `json:"name,omitempty"`
	Key string `json:"key,omitempty"`
}

// Process is an object t
//+hat specifies the post-processing or stop pre-processing.
type Process struct {
	// Process type. The value is command or http.
	// The command is to execute the command line, and http is to send an http request.
	Type string `json:"type" required:"true"`
	// Command parameters, such as ["sleep", "1"]. Applies to command type.
	Commands []string `json:"command,omitempty"`
	// The port number. Applies to http type.
	Port int `json:"port,omitempty"`
	// Request URL. Applies to http type.
	Path string `json:"path,omitempty"`
	// Defaults to the IP address of the POD instance. You can also specify it yourself. Applies to http type.
	Host string `json:"host,omitempty"`
}

// Affinity is an object that specifies the configuration details of the affinity or anti-affinity.
type Affinity struct {
	Kind string `json:"kind"`

	Condition string `json:"condition"`
	Weight int `json:"weight"`
	MatchExpressions []MatchExpression `json:"match_expressions"`
}

type MatchExpression struct {
	Key string `json:"key"`
	Value string `json:"value"`
	Operation string `json:"operation"`
}

type SecurityContext struct {
	RunAsUser int `json:"run_as_user"`
	RunAsGroup int `json:"run_as_group"`
	Capabilities Capability `json:"capabilities,omitempty"`
}

type Capability struct {
	Add []string `json:"add,omitempty"`
	Drop []string `json:"drop,omitempty"`
}

type DNSConfig struct {
	Searches []string `json:"searches,omitempty"`
	NameServers []string `json:"name_servers,omitempty"`
	Options []Variable `json:"options,omitempty"`
}

// Configuration is an object that specifies the environment variable for the component instance.
type Variable struct {
	// Environment variable name.
	// The value contains 1 to 64 characters, including letters, digits, underscores (_), hyphens (-), and dots (.),
	// and cannot start with a digit.
	Name string `json:"name" required:"true"`
	// Environment variable value.
	Value string `json:"value" required:"true"`
}

// LogCollectionPolicy is an object that specifies the policy of the log collection.
type LogCollectionPolicy struct {
	// Container mounting path.
	LogPath string `json:"logPath" required:"ture"`
	// Aging period.
	AgingPeriod string `json:"rotate" required:"ture"`
	// The extended host path, the valid values are as follows:
	//	None
	//	PodUID
	//	PodName
	//	PodUID/ContainerName
	//	PodName/ContainerName
	// If omited, means container mounting.
	HostExtendPath string `json:"hostExtendPath,omitempty"`
	// Host mounting path.
	HostPath string `json:"hostPath,omitempty"`
}

// ProbeDetail is an object that specifies the configuration details of the liveness probe and service probe.
type ProbeDetail struct {
	// Value: http, tcp, or command.
	// The check methods are HTTP request check, TCP port check, and command execution check, respectively.
	Type string `json:"type" required:"true"`
	// Interval between the startup and detection.
	Delay int `json:"delay,omitempty"`
	// Detection timeout interval.
	Timeout int `json:"timeout,omitempty"`
	Scheme string `json:"scheme,omitempty"`
	Host string `json:"host,omitempty"`
	Port int `json:"port,omitempty"`
	Path string `json:"path,omitempty"`
	PeriodSeconds int `json:"period_seconds,omitempty"`
	SuccessThreshold int `json:"success_threshold,omitempty"`
	FailureThreshold int `json:"failure_threshold,omitempty"`
	Command []string `json:"command,omitempty"`
	HttpHeaders []Variable `json:"http_headers,omitempty"`
}

type CustomMetric struct {
	Path string `json:"path,omitempty"`
	Port int `json:"port,omitempty"`
	Dimensions string `json:"dimensions,omitempty"`
}

// ReferResource is an object that specifies the deployed basic and optional resources.
type ReferResource struct {
	// Resource ID.
	// Note: If type is set to ecs, the value of this parameter must be Default.
	ID string `json:"id" required:"true"`
	// Resource type.
	// Basic resources: Cloud Container Engine (CCE), Auto Scaling (AS), and Elastic Cloud Server (ECS).
	// Optional resources: Relational Database Service (RDS), Distributed Cache Service (DCS),
	// Elastic Load Balance (ELB), and other services.
	Type string `json:"type" required:"true"`
	// Application alias, which is provided only in DCS scenario. Value: distributed_session, distributed_cache, or
	// distributed_session, distributed_cache. Default value: distributed_session, distributed_cache.
	ReferAlias string `json:"refer_alias,omitempty"`
	// Reference resource parameter.
	// NOTICE:
	// When type is set to cce, this parameter is mandatory. You need to specify the namespace of the cluster where the
	// component is to be deployed. Example: {"namespace": "default"}.
	// When type is set to ecs, this parameter is mandatory. You need to specify the hosts where the component is to be
	// deployed. Example: {"hosts":["04d9f887-9860-4029-91d1-7d3102903a69", "04d9f887-9860-4029-91d1-7d3102903a70"]}}.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type UpdateStrategy struct {
	Type string `json:"type,omitempty"`
	MaxUnavailable int `json:"max_unavailable,omitempty"`
	MaxSurge int `json:"max_surge,omitempty"`
}

type DeployStrategy struct {
	Type string `json:"type,omitempty"`
	RollingRelease RollingRelease `json:"rolling_release,omitempty"`
	GrayRelease GrayRelease `json:"gray_release,omitempty"`
}

type RollingRelease struct {
	Batches int `json:"batches,omitempty"`
}

type GrayRelease struct {
	Type string `json:"type,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Weight int `json:"weight,omitempty"`
	FirstBatchReplica int `json:"first_batch_replica,omitempty"`
	RemainingBatch int `json:"remaining_batch,omitempty"`
}

// Properties is an object to specified the other configuration of the software package for OBS bucket.
type Properties struct {
	// Object Storage Service (OBS) endpoint address. Example: https://obs.region_id.external_domain_name.com.
	Endpoint string `json:"endpoint,omitempty"`
	// Name of the OBS bucket where the software package is stored.
	Bucket string `json:"bucket,omitempty"`
	// Object in the OBS bucket, which is usually the name of the software package.
	// If there is a folder, the path of the folder must be added. Example: test.jar or demo/test.jar.
	Key string `json:"key,omitempty"`
}

// Parameter is an object to specified the building configuration of codes or package.
type Parameter struct {
}

// Create is a method to create a component in the specified appliation using given parameters.
func Create(c *golangsdk.ServiceClient, appId string, opts CreateOpts) (*JobResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Component
	_, err = c.Post(rootURL(c, appId), b, &r, nil)
	return &r, err
}

// Get is a method to retrieves a particular configuration based on its unique ID.
func Get(c *golangsdk.ServiceClient, appId, componentId string) (*Component, error) {
	var r Component
	_, err := c.Get(resourceURL(c, appId, componentId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// Value range: 0–100, or 1000.
	// Default value: 1000, indicating that a maximum of 1000 records can be queried and all records are displayed on
	// the same page.
	Limit string `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
	// Sorting field. By default, query results are sorted by creation time.
	// The following enumerated values are supported: create_time, name, and update_time.
	OrderBy string `q:"order_by"`
	// Descending or ascending order. Default value: desc.
	Order string `q:"order"`
}

// List is a method to query the list of the components using given opts.
func List(c *golangsdk.ServiceClient, appId string, opts ListOpts) ([]Component, error) {
	url := rootURL(c, appId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ComponentPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractComponents(pages)
}

// UpdateOpts is the structure required by the Update method to update the component configuration.
type UpdateOpts struct {
	// Application component name.
	// The value can contain 2 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	// It must start with a letter and end with a letter or digit.
	Name string `json:"name,omitempty"`
	// Description.
	// The value can contain up to 128 characters.
	Description *string `json:"description,omitempty"`
	// Source of the code or software package.
	Source *Source `json:"source,omitempty"`
	// Component build.
	Builder *Builder `json:"build,omitempty"`
}

// Update is a method to update the component configuration, such as name, description, builder and code source.
func Update(c *golangsdk.ServiceClient, appId, componentId string, opts UpdateOpts) (*Component, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Component
	_, err = c.Put(resourceURL(c, appId, componentId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing component from a specified application.
func Delete(c *golangsdk.ServiceClient, appId, componentId string) error {
	_, err := c.Delete(resourceURL(c, appId, componentId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
