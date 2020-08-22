# api
--
    import "."


## Usage

```go
const (
	Number  ServiceFieldType = "number"
	String                   = "string"
	Boolean                  = "boolean"
	Values                   = "values"
)
```

```go
var MessageHandlers = map[string]Message{}
```

```go
var ServiceNotFound = errors.New("not found")
```

#### func  RegisterMessageType

```go
func RegisterMessageType(typ ...Message)
```

#### type Action

```go
type Action interface {
	KeyFieldName() string
}
```


#### type AndCondition

```go
type AndCondition struct {
	Conditions []Condition `json:"conditions"`
}
```


#### func (*AndCondition) ConditionKey

```go
func (*AndCondition) ConditionKey() string
```

#### type AuthInvalidMessage

```go
type AuthInvalidMessage struct {
	Message string `json:"message"`
}
```


#### func (AuthInvalidMessage) Type

```go
func (AuthInvalidMessage) Type() string
```

#### type AuthMessage

```go
type AuthMessage struct {
	AccessToken string `json:"access_token,omitempty"`
	ApiPassword string `json:"api_password,omitempty"`
}
```


#### func (AuthMessage) Type

```go
func (AuthMessage) Type() string
```

#### type AuthOkMessage

```go
type AuthOkMessage struct{}
```


#### func (AuthOkMessage) Type

```go
func (AuthOkMessage) Type() string
```

#### type AuthRequiredMessage

```go
type AuthRequiredMessage struct{}
```


#### func (AuthRequiredMessage) Type

```go
func (AuthRequiredMessage) Type() string
```

#### type Automation

```go
type Automation struct {
	Action    AutomationAction      `json:"action"`
	Alias     string                `json:"alias"`
	Condition []AutomationCondition `json:"condition,omitempty"`
	Id        AutomationId          `json:"id"`
	Trigger   []AutomationTrigger   `json:"trigger"`
}
```


#### type AutomationAction

```go
type AutomationAction struct {
	Action
}
```


#### func (*AutomationAction) MarshalJSON

```go
func (a *AutomationAction) MarshalJSON() ([]byte, error)
```

#### func (*AutomationAction) UnmarshalJSON

```go
func (a *AutomationAction) UnmarshalJSON(data []byte) error
```

#### type AutomationCondition

```go
type AutomationCondition struct{ Condition }
```


#### func (AutomationCondition) MarshalJSON

```go
func (a AutomationCondition) MarshalJSON() ([]byte, error)
```

#### func (*AutomationCondition) UnmarshalJSON

```go
func (a *AutomationCondition) UnmarshalJSON(data []byte) error
```

#### type AutomationId

```go
type AutomationId string
```


#### type AutomationListEntry

```go
type AutomationListEntry struct {
	FriendlyName  string       `json:"friendly_name"`
	Id            AutomationId `json:"id"`
	LastTriggered time.Time    `json:"last_triggered"`
}
```


#### type AutomationTrigger

```go
type AutomationTrigger struct {
	Trigger
}
```


#### func (AutomationTrigger) MarshalJSON

```go
func (a AutomationTrigger) MarshalJSON() ([]byte, error)
```

#### func (*AutomationTrigger) UnmarshalJSON

```go
func (a *AutomationTrigger) UnmarshalJSON(data []byte) error
```

#### type Client

```go
type Client struct {
	Token  string
	Server string
}
```

Client is the object used to interface with a homeassistant server. It should
not be copied after being created.

#### func (*Client) Delete

```go
func (c *Client) Delete(path string, body interface{}) (interface{}, error)
```
Delete renders `body` as JSON and posts it to the given path.

#### func (*Client) Exchange

```go
func (c *Client) Exchange(send Message) (Message, error)
```
Exchange exchanges the 'send' Message for a response. It is safe for use by
multiple goroutines.

#### func (*Client) Get

```go
func (c *Client) Get(path string, parameters map[string]interface{}) (interface{}, error)
```

#### func (*Client) GetAutomation

```go
func (c *Client) GetAutomation(id AutomationId) (*Automation, error)
```

#### func (*Client) GetConfig

```go
func (c *Client) GetConfig() (*Config, error)
```

#### func (*Client) GetDevice

```go
func (c *Client) GetDevice(id string) (*Device, error)
```

#### func (*Client) GetEntity

```go
func (c *Client) GetEntity(id string) (*Entity, error)
```

#### func (*Client) GetFlow

```go
func (c *Client) GetFlow(id ConfigFlowId) (*ConfigFlow, error)
```
GetFlow returns the in-progress-but-not-started flow with the given ID.

#### func (*Client) GetOptionsFlow

```go
func (c *Client) GetOptionsFlow(id ConfigFlowId) (*ConfigFlow, error)
```
GetOptionsFlow gets the current status of the options flow with the given ID.

#### func (*Client) GetService

```go
func (c *Client) GetService(domain, service string) (*Service, error)
```
GetService returns the service with the given domain and service. If the service
can't be retrieved, the returned service will be nil and the error will be
non-nil. If the service does not exist but no other error occurs, an error
wrapping ServiceNotFound will be returned.

#### func (*Client) GetSystemOptions

```go
func (c *Client) GetSystemOptions(entryId EntryId) (*SystemOptions, error)
```

#### func (*Client) ListAutomations

```go
func (c *Client) ListAutomations() ([]AutomationListEntry, error)
```

#### func (*Client) ListConfigEntries

```go
func (c *Client) ListConfigEntries() ([]*ConfigEntry, error)
```
ListConfigEntries lists known config entries. A config entry is basically a
top-level device category; examples are `zwave` or `wemo`.

#### func (*Client) ListConfigFlowProgress

```go
func (c *Client) ListConfigFlowProgress() ([]ConfigFlowProgress, error)
```

#### func (*Client) ListDevices

```go
func (c *Client) ListDevices() ([]*Device, error)
```

#### func (*Client) ListEntities

```go
func (c *Client) ListEntities() ([]Entity, error)
```

#### func (*Client) ListFlowHandlers

```go
func (c *Client) ListFlowHandlers() ([]string, error)
```

#### func (*Client) ListServices

```go
func (c *Client) ListServices() ([]Service, error)
```

#### func (*Client) ListStates

```go
func (c *Client) ListStates() ([]State, error)
```

#### func (*Client) Post

```go
func (c *Client) Post(path string, body interface{}) (interface{}, error)
```
Post renders `body` as JSON and posts it to the given path.

#### func (*Client) Raw

```go
func (c *Client) Raw(method string, path string, parameters map[string]interface{}, body io.Reader) (interface{}, error)
```
Raw is as RawJSON, above, except the body is expected to already be an
io.Reader.

#### func (*Client) RawJSON

```go
func (c *Client) RawJSON(method string, path string, parameters map[string]interface{}, body interface{}) (interface{}, error)
```
RawJSON sends a request using the given method (e.g., GET, POST, DELETE) to the
given nominal path. Parameters describe URL parameters that are added to the
URL, and may be `nil`. body is an object that's converted to JSON and supplied
as the request body. No request body is supplied if `body` is nil.

Returns the generally JSON-decoded response object, or error, if something
happens. (Including, e.g., non-2XX responses or a body that isn't parseable
JSON).

#### func (*Client) RawRESTDelete

```go
func (c *Client) RawRESTDelete(path string, parameters map[string]interface{}) (interface{}, error)
```
RawRESTDelete requests the given path via the REST API and returns the
JSON-decoded request body if everything goes well. If anything else happens
(failure communicating, rejected request, etc.) the result will be nil and the
error will be non-nil. Note that the returned interface will be generically
typed (e.g., maps & slices). See RawRESTRequestAs.

#### func (*Client) RawRESTDeleteAs

```go
func (c *Client) RawRESTDeleteAs(path string, parameters map[string]interface{}, prototype interface{}) (interface{}, error)
```
RawRESTDeleteAs requests the given path via the REST API and returns the
JSON-decoded request body if everything goes well. If anything else happens
(failure communicating, rejected request, etc.) the result will be nil and the
error will be non-nil. The returned interface will be the same type as the given
prototype. If conversion cannot be achieved, an error will be returned.

#### func (*Client) RawRESTGet

```go
func (c *Client) RawRESTGet(path string, parameters map[string]interface{}) (interface{}, error)
```
RawRESTGet requests the given path via the REST API and returns the JSON-decoded
request body if everything goes well. If anything else happens (failure
communicating, rejected request, etc.) the result will be nil and the error will
be non-nil. Note that the returned interface will be generically typed (e.g.,
maps & slices). See RawRESTRequestAs.

#### func (*Client) RawRESTGetAs

```go
func (c *Client) RawRESTGetAs(path string, parameters map[string]interface{}, prototype interface{}) (interface{}, error)
```
RawRESTGetAs requests the given path via the REST API and returns the
JSON-decoded request body if everything goes well. If anything else happens
(failure communicating, rejected request, etc.) the result will be nil and the
error will be non-nil. The returned interface will be the same type as the given
prototype. If conversion cannot be achieved, an error will be returned.

#### func (*Client) RawRESTPost

```go
func (c *Client) RawRESTPost(path string, parameters map[string]interface{}) (interface{}, error)
```
RawRESTPost requests the given path via the REST API and returns the
JSON-decoded request body if everything goes well. If anything else happens
(failure communicating, rejected request, etc.) the result will be nil and the
error will be non-nil. Note that the returned interface will be generically
typed (e.g., maps & slices). See RawRESTRequestAs.

#### func (*Client) RawRESTPostAs

```go
func (c *Client) RawRESTPostAs(path string, parameters map[string]interface{}, prototype interface{}) (interface{}, error)
```
RawRESTPostAs requests the given path via the REST API and returns the
JSON-decoded request body if everything goes well. If anything else happens
(failure communicating, rejected request, etc.) the result will be nil and the
error will be non-nil. The returned interface will be the same type as the given
prototype. If conversion cannot be achieved, an error will be returned.

#### func (*Client) RawWebsocketRequest

```go
func (c *Client) RawWebsocketRequest(message Message) (interface{}, error)
```
RawWebsocketRequest exchanges the given message and returns the Result.Result if
everything goes well. If anything else happens (failure communicating, rejected
request, etc) the result will be nil and the error wll be non-nil. Note that the
returned interface will be generically typed (e.g., maps & slices). See
RawWebsocketRequestAs.

#### func (*Client) RawWebsocketRequestAs

```go
func (c *Client) RawWebsocketRequestAs(message Message, prototype interface{}) (interface{}, error)
```
RawWebsocketRequest exchanges the given message and returns the Result.Result if
everything goes well. If anything else happens (failure communicating, rejected
request, etc) the result will be nil and the error wll be non-nil. The returned
interface will be the same type as the given prototype. If conversion cannot be
achieved, an error will be returned.

#### func (*Client) SetEntityName

```go
func (c *Client) SetEntityName(id string, name string) error
```

#### func (*Client) SetFlow

```go
func (c *Client) SetFlow(id ConfigFlowId, payload map[string]interface{}) (result string, err error)
```
SetFlow sets the configuration for the given FlowId to the given payload, and
returns the result ID. I've observed result ID to refer to an Entry, but I guess
it could refer to a second-stage config flow?

If anything goes wrong, result will be the zero value and err will be non-nil.

#### func (*Client) StartOptionsFlow

```go
func (c *Client) StartOptionsFlow(entryId EntryId) (*ConfigFlow, error)
```
StartOptionsFlow initiates a options flow with the given handler, usually
(always?) a ConfigEntry id. The new flow is returned. The UI calls DELETE if a
thus-created flow is not used. It's not clear to me what happens if you don't do
this.

#### type Condition

```go
type Condition interface {
	ConditionKey() string
}
```


#### type Config

```go
type Config struct {
	Components            []string          `json:"components"`
	ConfigDir             string            `json:"config_dir"`
	ConfigSource          string            `json:"config_source"`
	Elevation             int               `json:"elevation"`
	Latitude              float64           `json:"latitude"`
	LocationName          string            `json:"location_name"`
	Longitude             float64           `json:"longitude"`
	TimeZone              string            `json:"time_zone"`
	UnitSystem            map[string]string `json:"unit_system"`
	Version               string            `json:"version"`
	WhitelistExternalDirs []string          `json:"whitelist_external_dirs"`
}
```


#### type ConfigEntry

```go
type ConfigEntry struct {
	EntryId         EntryId `json:"entry_id"`
	Domain          string  `json:"domain"`
	Title           string  `json:"title"`
	Source          string  `json:"source"`
	State           string  `json:"state"`
	ConnectionClass string  `json:"connection_class"`
	SupportsOptions bool    `json:"supports_options"`
}
```

https://github.com/home-assistant/core/blob/master/homeassistant/components/config/config_entries.py#L79

#### func (*ConfigEntry) GetSystemOptions

```go
func (c *ConfigEntry) GetSystemOptions() (*SystemOptions, error)
```

#### type ConfigFlow

```go
type ConfigFlow struct {
	DataSchema              []ConfigFlowDataSchema `json:"data_schema"`
	DescriptionPlaceholders interface{}            `json:"description_placeholders"`
	Errors                  map[string]string      `json:"errors"`
	FlowId                  ConfigFlowId           `json:"flow_id"`
	Handler                 string                 `json:"handler"`
	StepId                  string                 `json:"step_id"`
	Type                    string                 `json:"type"`

	// These fields appear in the response to SetFlow, maybe based on `Type`?
	Description string `json:"description,omitempty"`
	Result      string `json:"result,omitempty"`
	Title       string `json:"title,omitempty"`
}
```


#### type ConfigFlowDataSchema

```go
type ConfigFlowDataSchema struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
}
```


#### type ConfigFlowId

```go
type ConfigFlowId string
```


#### type ConfigFlowProgress

```go
type ConfigFlowProgress struct {
	Context ConfigFlowProgressContext `json:"context"`
	FlowId  ConfigFlowId              `json:"flow_id"`
	Handler string                    `json:"handler"`
}
```


#### type ConfigFlowProgressContext

```go
type ConfigFlowProgressContext struct {
	Source string `json:"source"`
}
```


#### type DelayAction

```go
type DelayAction struct {
	Delay int `json:"delay"`
}
```


#### func (*DelayAction) KeyFieldName

```go
func (*DelayAction) KeyFieldName() string
```

#### type Device

```go
type Device struct {
	ID            string      `json:"id"`
	AreaId        *string     `json:"area_id"`
	ConfigEntries []string    `json:"config_entries"`
	Connections   [][2]string `json:"connections"`
	Manufacturer  *string     `json:"manufacturer"`
	Model         *string     `json:"model"`
	Name          *string     `json:"name"`
	NameByUser    *string     `json:"name_by_user"`
	SwVersion     *string     `json:"sw_version"`
	ViaDeviceId   *string     `json:"via_device_id"`
}
```


#### type DeviceAction

```go
type DeviceAction struct {
	DeviceId string `json:"device_id"`
	Domain   string `json:"domain"`
	EntityId string `json:"entity_id"`
}
```


#### func (*DeviceAction) KeyFieldName

```go
func (*DeviceAction) KeyFieldName() string
```

#### type DeviceListMessage

```go
type DeviceListMessage struct{}
```


#### func (DeviceListMessage) Type

```go
func (DeviceListMessage) Type() string
```

#### type DeviceTrigger

```go
type DeviceTrigger struct {
	DeviceId string `json:"device_id"`
	Domain   string `json:"domain"`
	EntityId string `json:"entity_id"`
	Type     string `json:"type,omitempty"`
	Subtype  string `json:"subtype,omitempty"`
	Event    string `json:"event,omitempty"`
}
```


#### func (*DeviceTrigger) Platform

```go
func (*DeviceTrigger) Platform() string
```

#### type Entity

```go
type Entity struct {
	ConfigEntryId string `json:"config_entry_id,omitempty"`
	DeviceId      string `json:"device_id,omitempty"`
	DisabledBy    string `json:"disabled_by,omitempty"`
	EntityId      string `json:"entity_id,omitempty"`
	Platform      string `json:"platform,omitempty"`
	Name          string `json:"name,omitempty"`
}
```


#### type EntityGetMessage

```go
type EntityGetMessage struct {
	EntityId string `json:"entity_id"`
}
```


#### func (EntityGetMessage) Type

```go
func (EntityGetMessage) Type() string
```

#### type EntityList

```go
type EntityList struct{}
```


#### type EntityListMessage

```go
type EntityListMessage struct{}
```


#### func (EntityListMessage) Type

```go
func (EntityListMessage) Type() string
```

#### type EntityRename

```go
type EntityRename struct {
	EntityId string `json:"entity_id"`
	Name     string `json:"name"`
}
```


#### func (EntityRename) Type

```go
func (EntityRename) Type() string
```

#### type EntryId

```go
type EntryId string
```


#### type EventAction

```go
type EventAction struct {
	Event             string                 `json:"event"`
	EventData         map[string]interface{} `json:"event_data,omitempty"`
	EventDataTemplate map[string]interface{} `json:"event_data_template,omitempty"`
}
```


#### func (*EventAction) KeyFieldName

```go
func (*EventAction) KeyFieldName() string
```

#### type EventTrigger

```go
type EventTrigger struct {
	EventType string      `json:"event_type"`
	EventData interface{} `json:"event_data"`
}
```


#### func (*EventTrigger) Platform

```go
func (*EventTrigger) Platform() string
```

#### type GeoLocationTrigger

```go
type GeoLocationTrigger struct {
	Source string `json:"source"`
	Zone   string `json:"zone"`
	Event  string `json:"event"`
}
```


#### func (*GeoLocationTrigger) Platform

```go
func (*GeoLocationTrigger) Platform() string
```

#### type GetConfigMessage

```go
type GetConfigMessage struct{}
```


#### func (GetConfigMessage) Type

```go
func (GetConfigMessage) Type() string
```

#### type HassTrigger

```go
type HassTrigger struct {
	Event string `json:"event"`
}
```


#### func (*HassTrigger) Platform

```go
func (*HassTrigger) Platform() string
```

#### type ListConfigFlowProgressMessage

```go
type ListConfigFlowProgressMessage struct{}
```


#### func (ListConfigFlowProgressMessage) Type

```go
func (ListConfigFlowProgressMessage) Type() string
```

#### type ListStatesMessage

```go
type ListStatesMessage struct{}
```


#### func (ListStatesMessage) Type

```go
func (g ListStatesMessage) Type() string
```

#### type ListSystemOptionsMessage

```go
type ListSystemOptionsMessage struct {
	EntryId EntryId `json:"entry_id"`
}
```


#### func (ListSystemOptionsMessage) Type

```go
func (ListSystemOptionsMessage) Type() string
```

#### type Message

```go
type Message interface {
	Type() string
}
```


#### func  MessageFromJSON

```go
func MessageFromJSON(data []byte) (Message, error)
```
MessageFromJSON returns one of the various ___Message structs described above,
or an error if we can't do so.

#### type MqttTrigger

```go
type MqttTrigger struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload,omitempty"`
}
```


#### func (*MqttTrigger) Platform

```go
func (*MqttTrigger) Platform() string
```

#### type NotCondition

```go
type NotCondition struct {
	Conditions []Condition `json:"conditions"`
}
```


#### func (*NotCondition) ConditionKey

```go
func (*NotCondition) ConditionKey() string
```

#### type NumericStateCondition

```go
type NumericStateCondition struct {
	EntityId string  `json:"entity_id"`
	Above    float64 `json:"above"`
}
```


#### func (*NumericStateCondition) ConditionKey

```go
func (*NumericStateCondition) ConditionKey() string
```

#### type NumericStateTrigger

```go
type NumericStateTrigger struct {
	EntityId      string        `json:"entity_id"`
	Above         float64       `json:"above"`
	Below         float64       `json:"below"`
	ValueTemplate string        `json:"value_template"`
	For           time.Duration `json:"for"`
}
```


#### func (*NumericStateTrigger) Platform

```go
func (*NumericStateTrigger) Platform() string
```

#### type OrCondition

```go
type OrCondition struct {
	Conditions []Condition `json:"conditions"`
}
```


#### func (*OrCondition) ConditionKey

```go
func (*OrCondition) ConditionKey() string
```

#### type ResultError

```go
type ResultError struct {
	Code    string
	Message string
}
```


#### type ResultMessage

```go
type ResultMessage struct {
	Id      int
	Success bool
	Result  interface{}
	Error   ResultError `json:"error,omitempty"`
}
```


#### func (ResultMessage) Type

```go
func (ResultMessage) Type() string
```

#### type SceneAction

```go
type SceneAction struct {
	Scene string `json:"scene"`
}
```


#### func (*SceneAction) KeyFieldName

```go
func (*SceneAction) KeyFieldName() string
```

#### type Service

```go
type Service struct {
	Domain      string                   `json:"domain"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Fields      map[string]*ServiceField `json:"fields"`
}
```


#### func (*Service) Call

```go
func (s *Service) Call(data map[string]interface{}) ([]State, error)
```

#### type ServiceAction

```go
type ServiceAction struct {
	Service  string                 `json:"service"`
	EntityId *string                `json:"entity_id,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}
```


#### func (*ServiceAction) KeyFieldName

```go
func (*ServiceAction) KeyFieldName() string
```

#### type ServiceField

```go
type ServiceField struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        ServiceFieldType `json:"type"`
	Example     interface{}      `json:"example"`
	Values      []interface{}    `json:"values"`
	Default     interface{}      `json:"default"`
}
```


#### func (*ServiceField) UnmarshalJSON

```go
func (s *ServiceField) UnmarshalJSON(data []byte) error
```

#### type ServiceFieldType

```go
type ServiceFieldType string
```


#### func (ServiceFieldType) String

```go
func (s ServiceFieldType) String() string
```

#### type State

```go
type State struct {
	Attributes map[string]interface{} `json:"attributes"`
	Context    struct {
		Id       string  `json:"id"`
		ParentId *string `json:"parent_id"`
		UserId   *string `json:"user_id"`
	} `json:"context"`
	EntityId    string    `json:"entity_id"`
	LastChanged time.Time `json:"last_Changed"`
	LastUpdated time.Time `json:"last_updated"`
	State       string    `json:"state"`
}
```


#### type StateCondition

```go
type StateCondition struct {
	EntityId string        `json:"entity_id"`
	State    StringOrFloat `json:"state"`
}
```


#### func (*StateCondition) ConditionKey

```go
func (*StateCondition) ConditionKey() string
```

#### type StateTrigger

```go
type StateTrigger struct {
	EntityId string        `json:"entity_id,omitempty"`
	From     StringOrFloat `json:"from"`
	To       StringOrFloat `json:"to"`
	For      time.Duration `json:"for"`
}
```


#### func (*StateTrigger) Platform

```go
func (*StateTrigger) Platform() string
```

#### type StringOrFloat

```go
type StringOrFloat struct {
}
```

StringOrFloat represents a value that is either a string or an integer.

#### func  StringOrIntFromFloat64

```go
func StringOrIntFromFloat64(f float64) StringOrFloat
```

#### func  StringOrIntFromString

```go
func StringOrIntFromString(s string) StringOrFloat
```

#### func (StringOrFloat) MarshalJSON

```go
func (s StringOrFloat) MarshalJSON() ([]byte, error)
```

#### func (StringOrFloat) String

```go
func (s StringOrFloat) String() string
```

#### func (StringOrFloat) UnmarshalJSON

```go
func (s StringOrFloat) UnmarshalJSON(data []byte) error
```

#### type SunCondition

```go
type SunCondition struct {
	AfterOffset  float64 `json:"after_offset"`
	BeforeOffset float64 `json:"before_offset"`
	After        string  `json:"after"`
	Before       string  `json:"before"`
}
```


#### func (*SunCondition) ConditionKey

```go
func (*SunCondition) ConditionKey() string
```

#### type SunTrigger

```go
type SunTrigger struct {
	Offset float64 `json:"offset"`
	Event  string  `json:"event"`
}
```


#### func (*SunTrigger) Platform

```go
func (*SunTrigger) Platform() string
```

#### type SystemOptions

```go
type SystemOptions struct {
	DisableNewEntities bool `json:"disable_new_entities"`
}
```


#### type TemplateCondition

```go
type TemplateCondition struct {
	ValueTemplate string `json:"value_template"`
}
```


#### func (*TemplateCondition) ConditionKey

```go
func (*TemplateCondition) ConditionKey() string
```

#### type TemplateTrigger

```go
type TemplateTrigger struct {
	ValueTemplate string `json:"value_template"`
}
```


#### func (*TemplateTrigger) Platform

```go
func (*TemplateTrigger) Platform() string
```

#### type TimeCondition

```go
type TimeCondition struct {
	After  string `json:"after"`
	Before string `json:"before"`
}
```


#### func (*TimeCondition) ConditionKey

```go
func (*TimeCondition) ConditionKey() string
```

#### type TimePatternTrigger

```go
type TimePatternTrigger struct {
	Hours   StringOrFloat `json:"hours"`
	Minutes StringOrFloat `json:"minutes"`
	Seconds StringOrFloat `json:"seconds"`
}
```


#### func (*TimePatternTrigger) Platform

```go
func (*TimePatternTrigger) Platform() string
```

#### type TimeTrigger

```go
type TimeTrigger struct {
	// A time like HH:MM:SS, 24-hour time.
	At string `json:"at"`
}
```


#### func (*TimeTrigger) Platform

```go
func (*TimeTrigger) Platform() string
```

#### type Trigger

```go
type Trigger interface {
	Platform() string
}
```


#### type WaitAction

```go
type WaitAction struct {
	WaitTemplate string `json:"wait_template"`
	Timeout      int    `json:"timeout,omitempty"`
}
```


#### func (*WaitAction) KeyFieldName

```go
func (*WaitAction) KeyFieldName() string
```

#### type WebhookTrigger

```go
type WebhookTrigger struct {
	WebhookId string `json:"webhook_id"`
}
```


#### func (*WebhookTrigger) Platform

```go
func (*WebhookTrigger) Platform() string
```

#### type ZoneCondition

```go
type ZoneCondition struct {
	EntityId string `json:"entity_id"`
	Zone     string `json:"zone"`
}
```


#### func (*ZoneCondition) ConditionKey

```go
func (*ZoneCondition) ConditionKey() string
```

#### type ZoneTrigger

```go
type ZoneTrigger struct {
	EntityId string `json:"entity_id"`
	Zone     string `json:"zone"`
	Event    string `json:"event"`
}
```


#### func (*ZoneTrigger) Platform

```go
func (*ZoneTrigger) Platform() string
```
