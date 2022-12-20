package servicestage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chnsz/golangsdk/openstack/servicestage/v2/instances"
	"github.com/chnsz/golangsdk/openstack/servicestage/v3/components"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"log"
	"regexp"
	"time"
)

func lifecycleProcessSchemaResource1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"command", "http",
				}, false),
			},
			"commands": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func affinitySchemaResource1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"match_expression": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func probeDetailSchemaResource1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"command", "http", "tcp",
				}, false),
			},
			"command_param": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commands": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"scheme": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP", "HTTPS",
				}, false),
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"success_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"failure_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"http_header": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func ResourceComp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCompCreate,
		ReadContext:   resourceCompRead,
		UpdateContext: resourceCompUpdate,
		DeleteContext: resourceCompDelete,
		//
		//Importer: &schema.ResourceImporter{
		//	StateContext: resourceComponentImportState,
		//},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]([\w-]*[A-Za-z0-9])?$`),
						"The name can only contain letters, digits, underscores (_) and hyphens (-), and the name must"+
							" start with a letter and end with a letter or digit."),
					validation.StringLenBetween(2, 64),
				),
			},
			"source": {
				Type:     schema.TypeList,
				Required: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"code", "package", "image",
							}, false),
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GitHub", "GitLab", "Gitee", "Bitbucket", "package", "DevCloud",
							}, false),
						},
						"repo_url": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"authorization": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ExactlyOneOf: []string{"source.0.storage_type"},
						},
						"repo_ref": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"source.0.storage_type"},
						},
						"repo_namespace": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"source.0.storage_type"},
						},
					},
				},
			},
			"runtime": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"deploy_mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"container",
							}, false),
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workload_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"jvm_opts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tomcat_opts": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_xml": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"builder": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization": {
							Type:     schema.TypeString,
							Required: true,
						},
						"environment_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dockerfile_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cmd": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"use_public_cluster": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"node_label": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
										Computed: true,
									},
									"valye": {
										Type:     schema.TypeString,
										Required: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"workload_kind": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit_cpu": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"limit_memory": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"request_cpu": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"request_memory": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"replica": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ports": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"target_port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"external_access": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"acceleration": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"claim_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"env_variable": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_.]([\w-.]*)?$`),
												"The name can only contain letters, digits, underscores (_), "+
													"hyphens (-) and dots (.), and cannot start with a digit."),
											validation.StringLenBetween(1, 64),
										),
									},
									"inner": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value_form": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"reference_type": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"key": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"start_command": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"commands": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"args": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"storage": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"HostPath", "EmptyDir", "ConfigMap", "Secret", "PersistentVolumeClaim",
										}, false),
									},
									"parameter": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"claim_name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"secret_name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"medium": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"mount": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
												"readonly": {
													Type:     schema.TypeBool,
													Required: true,
												},
												"sub_path": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"lifecycle": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"post_start": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     lifecycleProcessSchemaResource1(),
									},
									"pre_stop": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     lifecycleProcessSchemaResource1(),
									},
								},
							},
						},
						"scheduler": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"affinity": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     affinitySchemaResource1(),
									},
									"anti_affinity": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     affinitySchemaResource1(),
									},
								},
							},
						},
						"security_context": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"run_as_user": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"run_as_group": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"capability": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"add": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"drop": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"dns_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_policy": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"searches": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 3,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"name_servers": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 2,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"options": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},

						"log_collection_policy": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_mounting": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
												"host_extend_path": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"aging_period": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "Hourly",
													ValidateFunc: validation.StringInSlice([]string{
														"Hourly", "Daily", "Weekly",
													}, false),
												},
											},
										},
									},
									"host_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"probe": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"liveness": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     probeDetailSchemaResource1(),
									},
									"readiness": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     probeDetailSchemaResource1(),
									},
								},
							},
						},
						"custom_metric": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"dimensions": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},

						"labels": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"strategy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"upgrade": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "RollingUpdate",
													ValidateFunc: validation.StringInSlice([]string{
														"RollingUpdate", "Recreate",
													}, false),
												},
												"max_unavailable": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
												"max_surge": {
													Type:     schema.TypeInt,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"deploy": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "RollingUpdate",
													ValidateFunc: validation.StringInSlice([]string{
														"OneBatchRelease", "RollingRelease", "GrayRelease",
													}, false),
												},
												"rolling_release": {
													Type:     schema.TypeSet,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"batches": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
												"gray_release": {
													Type:     schema.TypeSet,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"headers": {
																Type:     schema.TypeString,
																Optional: true,
																Computed: true,
															},
															"weight": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
															"first_batch_replica": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
															"remaining_batch": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"refer_resource": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceCompCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ServiceStageV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	opt, err := buildCompCreateOpts(d)
	if err != nil {
		return diag.Errorf("error building the CreateOpts of the component instance: %s", err)
	}
	log.Printf("[DEBUG] The instance create option of ServiceStage component is: %v", opt)

	resp, err := components.Create(client, appId, opt)
	if err != nil {
		return diag.Errorf("error creating ServiceStage component instance: %s", err)
	}

	d.SetId(resp.ComponentId)

	log.Printf("[DEBUG] Waiting for the component instance to become running, the instance ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      componentInstanceRefreshFunc(client, resp.JobId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the creation of component instance (%s) to complete: %s",
			d.Id(), err)
	}

	return resourceCompRead(ctx, d, meta)
}

func resourceCompRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	componentId := d.Get("component_id").(string)
	resp, err := instances.Get(client, appId, componentId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage component instance")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("environment_id", resp.EnvironmentId),
		d.Set("name", resp.Name),
		d.Set("version", resp.Version),
		d.Set("replica", resp.StatusDetail.Replica),
		d.Set("flavor_id", resp.FlavorId),
		d.Set("description", resp.Description),
		d.Set("artifact", flattenArtifact(resp.Artifacts)),
		d.Set("refer_resource", flattenReferResources(resp.ReferResources)),
		d.Set("configuration", flattenConfiguration(resp.Configuration)),
		d.Set("external_access", flattenExternalAccesses(resp.ExternalAccesses)),
		// Attributes
		d.Set("status", resp.StatusDetail.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}


func resourceCompUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	componentId := d.Get("component_id").(string)
	opt, err := buildInstanceUpdateOpts(d)
	if err != nil {
		return diag.Errorf("error building the UpdateOpts of the component instance: %s", err)
	}
	log.Printf("[DEBUG] The instance update option of ServiceStage component is: %v", opt)

	resp, err := instances.Update(client, appId, componentId, d.Id(), opt)
	if err != nil {
		return diag.Errorf("error updating component instance: %s", err)
	}

	log.Printf("[DEBUG] Waiting for the component instance to become running, the instance ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      componentInstanceRefreshFunc(client, resp.JobId),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the updation of component instance (%s) to complete: %s",
			d.Id(), err)
	}

	return resourceComponentInstanceRead(ctx, d, meta)
}

func resourceCompDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	componentId := d.Get("component_id").(string)
	resp, err := instances.Delete(client, appId, componentId, d.Id())
	if err != nil {
		return diag.Errorf("error deleting ServiceStage component instance (%s): %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Waiting for the component instance to become deleted, the instance ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      componentInstanceRefreshFunc(client, resp.JobId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the delete of component instance (%s) to complete: %s",
			d.Id(), err)
	}

	return nil
}


func buildCompCreateOpts(d *schema.ResourceData) (components.CreateOpts, error) {

	configuration := buildConfigurationStructure1(d.Get("configuration").(*schema.Set))

	resourceConfiguration := buildResourceConfigurationStructure(d.Get("resource_configuration").(*schema.Set))
	referResources,err:=buildReferResourcesList1(d.Get("refer_resources").(*schema.Set))
	if err != nil {
		return components.CreateOpts{}, err
	}
	result := components.CreateOpts{
		EnterpriseProjectId: "",
		Name:                d.Get("name").(string),
		Version:             d.Get("version").(string),
		EnvironmentId:       d.Get("environment_id").(string),
		WorkloadName:        d.Get("workload_name").(string),
		ApplicationId:       d.Get("application_id").(string),
		Description:         d.Get("description").(string),
		JVMOpts:             d.Get("jvm_opts").(string),
		RuntimeStack:        buildRuntimeStructure(d.Get("runtime_stack").([]interface{})),
		TomcatOpts:          buildTomcatOptsStructure(d.Get("tomcat_opts").(*schema.Set)),
		Builder:             buildBuilderStructure(d.Get("builder").([]interface{})),
		Source:              buildSourceStructure(d.Get("source").([]interface{})),
		WorkloadKind:        d.Get("workload_kind").(string),
		LimitCPU:            resourceConfiguration["limit_cpu"].(int),
		LimitMemory:         resourceConfiguration["limit_memory"].(int),
		RequestCPU:          resourceConfiguration["request_cpu"].(int),
		RequestMemory:       resourceConfiguration["request_memory"].(int),
		Replica:             resourceConfiguration["replica"].(int),
		ServiceName:         d.Get("service_name").(string),

		Ports:            buildPortsStructure(d.Get("port").([]interface{})),
		ExternalAccesses: buildExternalAccessStructure(d.Get("external_access").(*schema.Set)),
		Acceleration:     buildAccelerationStructure(d.Get("acceleration").(*schema.Set)),
		EnvVariables:     configuration.EnvVariables,
		Commands:         configuration.Commands,
		Storages:        configuration.Storages,
		PostStart:        configuration.PostStart,
		PreStop:          configuration.PreStop,
		Affinity:         configuration.Affinity,
		AntiAffinity:     configuration.AntiAffinity,
		SecurityContext:  configuration.SecurityContext,

		DNSPolicy: configuration.DNSPolicy,

		DNSConfig:             configuration.DNSConfig,
		LogCollectionPolicies: configuration.LogCollectionPolicies,
		LivenessProbe:         configuration.LivenessProbe,
		ReadinessProbe:        configuration.ReadinessProbe,
		CustomMetric:          configuration.CustomMetric,
		ReferResources:        referResources,
		//Labels []string `json:"labels,omitempty"`
		//UpdateStrategy UpdateStrategy `json:"update_strategy,omitempty"`
		//DeployStrategy DeployStrategy `json:"deploy_strategy,omitempty"`
	}
	return result, nil
}

func buildRuntimeStructure(params []interface{}) *components.RuntimeStack {
	if len(params) == 0 {
		return nil
	}
	val := params[0].(map[string]interface{})

	return &components.RuntimeStack{
		Name:       val["name"].(string),
		Version:    val["version"].(string),
		DeployMode: val["deploy_mode"].(string),
		Type:       val["type"].(string),
	}
}

func buildTomcatOptsStructure(params *schema.Set) *components.TomcatOptions {
	if params.Len() < 1 {
		return nil
	}
	value := params.List()[0].(map[string]interface{})
	return &components.TomcatOptions{
		ServerXML: value["server_xml"].(string),
	}
}

func buildBuilderStructure(params []interface{}) *components.Builder {
	if len(params) < 1 {
		return nil
	}

	param := params[0].(map[string]interface{})
	return &components.Builder{
		BuildCmd:          param["cmd"].(string),
		DockerfilePath:    param["dockerfile_path"].(string),
		ArtifactNamespace: param["organization"].(string),
		UsePublicCluster:  param["use_public_cluster"].(bool),
		EnvironmentId:     param["environment_id"].(string),
		NodeLabelSelector: buildNodeLabelStructure(param["node_label"].(*schema.Set)),
	}
}

func buildNodeLabelStructure(params *schema.Set) map[string]interface{} {
	if params.Len() < 1 {
		return nil
	}
	nodeLabel := params.List()[0].(map[string]interface{})
	return map[string]interface{}{
		nodeLabel["key"].(string): nodeLabel["value"],
	}
}

func buildSourceStructure(params []interface{}) *components.Source {
	if len(params) < 1 {
		return nil
	}
	value := params[0].(map[string]interface{})
	return &components.Source{
		Kind:          value["kind"].(string),
		Url:           value["url"].(string),
		Version:       value["version"].(string),
		Storage:       value["storage"].(string),
		RepoType:      value["type"].(string),
		RepoUrl:       value["repo_url"].(string),
		RepoAuth:      value["authorization"].(string),
		RepoNamespace: value["repo_namespace"].(string),
		RepoRef:       value["repo_ref"].(string),
	}
}

func buildResourceConfigurationStructure(params *schema.Set) map[string]interface{} {
	if params.Len() < 1 {
		return nil
	}
	val := params.List()[0].(map[string]interface{})
	if _, ok := val["limit_cpu"]; !ok {
		val["limit_cpu"] = 0.25
	}
	if _, ok := val["limit_memory"]; !ok {
		val["limit_memory"] = 0.25
	}
	if _, ok := val["request_cpu"]; !ok {
		val["request_cpu"] = 0.25
	}
	if _, ok := val["request_memory"]; !ok {
		val["request_memory"] = 0.25
	}
	if _, ok := val["replica"]; !ok {
		val["replica"] = 1
	}
	return val
}

func buildPortsStructure(params []interface{}) []components.Port {
	if len(params) < 1 {
		return nil
	}
	result := make([]components.Port, 0, len(params))
	for _, value := range params {
		val := value.(map[string]interface{})
		result = append(result, components.Port{
			Name:       val["name"].(string),
			Port:       val["port"].(int),
			TargetPort: val["target_port"].(int),
		})
	}
	return result
}

func buildExternalAccessStructure(params *schema.Set) []components.ExternalAccess {
	if params.Len() < 1 {
		return nil
	}

	result := make([]components.ExternalAccess, params.Len())
	for i, val := range params.List() {
		access := val.(map[string]interface{})
		result[i] = components.ExternalAccess{
			Protocol:    access["protocol"].(string),
			Address:     access["address"].(string),
			ForwardPort: access["port"].(int),
			ID:          access["id"].(string),
			Type:        access["type"].(string),
		}
	}

	return result
}

func buildAccelerationStructure(params *schema.Set) *components.Acceleration {
	if params.Len() < 1 {
		return nil
	}
	value := params.List()[0].(map[string]interface{})
	return &components.Acceleration{
		ClaimName: value["claim_name"].(string),
	}
}

func buildConfigurationStructure1(params *schema.Set) *Configuration {
	if params.Len() < 1 {
		return nil
	}
	value := params.List()[0].(map[string]interface{})

	lifecycle := buildLifecycleStructure1(value["lifecycle"].(*schema.Set))
	scheduler := buildSchedulerStructure1(value["scheduler"].(*schema.Set))
	dnsPolicy, dnsConfig := buildDNSConfigStructure(value["dns_config"].(*schema.Set))
	probe := buildProbeStructure1(value["probe"].([]interface{}))

	return &Configuration{
		EnvVariables:    buildEnvVariablesStructure(value["env"].([]interface{})),
		Commands:        buildCommandStructure(value["start_command"].(*schema.Set)),
		Storages:        buildStorageStructure(value["storage"].([]interface{})),
		PostStart:       lifecycle["post_start"],
		PreStop:         lifecycle["pre_stop"],
		Affinity:        scheduler["affinity"],
		AntiAffinity:    scheduler["anti-affinity"],
		SecurityContext: buildSecurityContextStructure(value["security_context"].(*schema.Set)),

		DNSPolicy: dnsPolicy,

		DNSConfig:             dnsConfig,
		LogCollectionPolicies: buildLogCollectionPoliciesStructure1(value["log"].(*schema.Set)),
		LivenessProbe :probe["liveness"],
		ReadinessProbe:probe["readiness"],
		CustomMetric:buildCustomMetricStructure(value["custom_metric"].(*schema.Set)),
	}

}

func buildEnvVariablesStructure(params []interface{}) []components.EnvVariable {
	if len(params) < 1 {
		return nil
	}

	result := make([]components.EnvVariable, 0, len(params))
	for _, value := range params {
		val := value.(map[string]interface{})
		result = append(result, components.EnvVariable{
			Name:      val["name"].(string),
			Inner:     val["inner"].(bool),
			Value:     val["value"].(string),
			ValueFrom: buildValueFromStructure(val["value_from"].(map[string]interface{})),
		})
	}
	return result
}

func buildValueFromStructure(params map[string]interface{}) components.ValueFrom {
	valueFrom := components.ValueFrom{}
	if val, ok := params["reference_type"]; ok {
		valueFrom.ReferenceType = val.(string)
	}
	if val, ok := params["name"]; ok {
		valueFrom.Name = val.(string)
	}
	if val, ok := params["key"]; ok {
		valueFrom.Key = val.(string)
	}
	return valueFrom
}

func buildCommandStructure(params *schema.Set) *components.Command {
	if params.Len() < 1 {
		return nil
	}
	value := params.List()[0].(map[string]interface{})

	return &components.Command{
		Command: value["command"].(string),
		Args:    value["args"].([]string),
	}
}

func buildStorageStructure(params []interface{}) []components.Storage {
	if len(params) < 1 {
		return nil
	}
	result := make([]components.Storage, 0, len(params))
	for _, value := range params {
		val := value.(map[string]interface{})
		result = append(result, components.Storage{
			Type:       val["type"].(string),
			Parameters: buildStorageParameterStructure(val["parameter"].(*schema.Set)),
			Mounts:     buildStorageMountStructure(val["mounts"].(*schema.Set)),
		})
	}
	return result
}

func buildStorageParameterStructure(params *schema.Set) *components.StorageParams {
	if params.Len() < 1 {
		return nil
	}
	value := params.List()[0].(map[string]interface{})
	return &components.StorageParams{
		Path:        value["path"].(string),
		Name:        value["name"].(string),
		DefaultMode: value["default_mode"].(string),
		Medium:      value["medium"].(string),
	}
}

func buildStorageMountStructure(params *schema.Set) *components.Mount {
	if params.Len() < 1 {
		return nil
	}
	value := params.List()[0].(map[string]interface{})
	return &components.Mount{
		Path:     value["path"].(string),
		SubPath:  value["sub_path"].(string),
		Readonly: value["readonly"].(bool),
	}
}

func buildLifecycleStructure1(params *schema.Set) map[string]*components.Process {
	if params.Len() < 1 {
		return nil
	}
	result := make(map[string]*components.Process)
	value := params.List()[0].(map[string]interface{})

	for k, v := range value {
		val := v.(map[string]interface{})
		result[k] = &components.Process{
			Type:     val["type"].(string),
			Host:     val["host"].(string),
			Port:     val["port"].(int),
			Path:     val["path"].(string),
			Commands: val["command"].([]string),
		}
	}
	return result
}

func buildSchedulerStructure1(params *schema.Set) map[string][]components.Affinity {
	if params.Len() < 1 {
		return nil
	}
	result := make(map[string][]components.Affinity)
	value := params.List()[0].(map[string]interface{})

	for k, v := range value {
		result[k] = buildAffinityStructure1(v.(*schema.Set))
	}
	return result
}

func buildAffinityStructure1(params *schema.Set) []components.Affinity {
	if params.Len() < 1 {
		return nil
	}

	result := make([]components.Affinity, 0, params.Len())
	for _, v := range params.List() {
		val := v.(map[string]interface{})
		result = append(result, components.Affinity{
			Kind:             val["kind"].(string),
			Condition:        val["condition"].(string),
			Weight:           val["weight"].(int),
			MatchExpressions: buildMatchExpressionStructure(val["match_expression"].(*schema.Set)),
		})
	}
	return result
}

func buildMatchExpressionStructure(params *schema.Set) []components.MatchExpression {
	if params.Len() < 1 {
		return nil
	}
	result := make([]components.MatchExpression, 0, params.Len())
	for _, v := range params.List() {
		val := v.(map[string]interface{})
		result = append(result, components.MatchExpression{
			Key:       val["key"].(string),
			Value:     val["value"].(string),
			Operation: val["operation"].(string),
		})
	}
	return result
}

func buildSecurityContextStructure(params *schema.Set) *components.SecurityContext {
	if params.Len() < 1 {
		return nil
	}
	value := params.List()[0].(map[string]interface{})
	return &components.SecurityContext{
		RunAsUser:    value["run_as_user"].(int),
		RunAsGroup:   value["run_as_group"].(int),
		Capabilities: buildCapabilityStructure(value["capability"].(*schema.Set)),
	}
}

func buildCapabilityStructure(params *schema.Set) components.Capability {
	if params.Len() < 1 {
		return components.Capability{}
	}
	value := params.List()[0].(map[string]interface{})
	return components.Capability{
		Add:  value["add"].([]string),
		Drop: value["drop"].([]string),
	}
}

func buildDNSConfigStructure(params *schema.Set) (string, *components.DNSConfig) {
	if params.Len() < 1 {
		return "", nil
	}
	value := params.List()[0].(map[string]interface{})
	return value["dns_policy"].(string), &components.DNSConfig{
		Searches:    value["searches"].([]string),
		NameServers: value["name_servers"].([]string),
		Options:     buildEnvVariables1(value["options"].(*schema.Set)),
	}
}

func buildEnvVariables1(params *schema.Set) []components.Variable {
	if params.Len() < 1 {
		return nil
	}

	result := make([]components.Variable, 0, params.Len())
	for _, val := range params.List() {
		variable := val.(map[string]interface{})
		result = append(result, components.Variable{
			Name:  variable["name"].(string),
			Value: variable["value"].(string),
		})
	}

	return result
}

func buildLogCollectionPoliciesStructure1(params *schema.Set) []components.LogCollectionPolicy {
	if params.Len() < 1 {
		return nil
	}

	result := make([]components.LogCollectionPolicy, 0, params.Len())
	for _, val := range params.List() {
		policy := val.(map[string]interface{})
		hostPath := policy["host_path"].(string)
		cmSet := policy["container_mounting"].(*schema.Set)
		for _, val := range cmSet.List() {
			cm := val.(map[string]interface{})
			result = append(result, components.LogCollectionPolicy{
				LogPath:        cm["path"].(string),
				HostExtendPath: cm["host_extend_path"].(string),
				AgingPeriod:    cm["aging_period"].(string),
				HostPath:       hostPath,
			})
		}
	}

	return result
}

func buildProbeStructure1(params []interface{}) map[string]*components.ProbeDetail {
	if len(params) < 1 {
		return nil
	}
	result := make(map[string]*components.ProbeDetail)
	probe := params[0].(map[string]interface{})
	result["liveness"] = buildProbeDetailStructure1(probe["liveness"].([]interface{}))
	result["readiness"] = buildProbeDetailStructure1(probe["readiness"].([]interface{}))
	return result
}

func buildProbeDetailStructure1(params []interface{}) *components.ProbeDetail {
	if len(params) < 1 {
		return nil
	}
	value := params[0].(map[string]interface{})
	return &components.ProbeDetail{
		Type:             value["type"].(string),
		Delay:            value["delay"].(int),
		Timeout:          value["timeout"].(int),
		Scheme:           value["scheme"].(string),
		Host:             value["host"].(string),
		Port:             value["port"].(int),
		Path:             value["path"].(string),
		PeriodSeconds:    value["period_seconds"].(int),
		SuccessThreshold: value["success_threshold"].(int),
		FailureThreshold: value["failure_threshold"].(int),
		Command:          value["command"].([]string),
		HttpHeaders:      buildEnvVariables1(value["http_header"].(*schema.Set)),
	}
}

func buildCustomMetricStructure(params *schema.Set)*components.CustomMetric{
	if params.Len()<1{
		return nil
	}
	value:=params.List()[0].(map[string]interface{})
	return &components.CustomMetric{
		Path: value["path"].(string),
		Port: value["port"].(int),
		Dimensions: value["dimensions"].(string),
	}
}

type Configuration struct {
	EnvVariables    []components.EnvVariable
	Commands        *components.Command
	Storages        []components.Storage
	PostStart       *components.Process
	PreStop         *components.Process
	Affinity        []components.Affinity
	AntiAffinity    []components.Affinity
	SecurityContext *components.SecurityContext

	DNSPolicy string

	DNSConfig             *components.DNSConfig
	LogCollectionPolicies []components.LogCollectionPolicy
	LivenessProbe         *components.ProbeDetail
	ReadinessProbe        *components.ProbeDetail
	CustomMetric          *components.CustomMetric
	ReferResources        *schema.Set
}

func buildReferResourcesList1(params *schema.Set) ([]components.ReferResource, error) {
	if params.Len() < 1 {
		return nil, nil
	}

	result := make([]components.ReferResource, params.Len())
	for i, val := range params.List() {
		res := val.(map[string]interface{})
		refer := components.ReferResource{
			Type:       res["type"].(string),
			ID:         res["id"].(string),
			ReferAlias: res["alias"].(string),
		}
		pResult := make(map[string]interface{})
		if param, ok := res["parameters"]; ok {
			log.Printf("[DEBUG] The parameters is %#v", param)
			p := param.(map[string]interface{})
			for k, v := range p {
				if k == "hosts" {
					var r []string
					err := json.Unmarshal([]byte(v.(string)), &r)
					if err != nil {
						return nil, fmt.Errorf("the format of the host value is not right: %#v", v)
					}
					pResult[k] = &r
					continue
				}
				pResult[k] = v
			}

			refer.Parameters = pResult
		}

		log.Printf("[DEBUG] The parameter map is %v", pResult)
		result[i] = refer
	}

	return result, nil
}
