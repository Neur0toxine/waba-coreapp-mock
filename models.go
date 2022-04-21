package main

import "encoding/json"

type (
	BlockingType             string
	ContactStatus            string
	MessageType              string
	RecipientType            string
	ContactAddressType       string
	ContactEmailType         string
	ContactPhoneType         string
	ContactURLType           string
	TemplateLanguagePolicy   string
	TemplateComponentType    string
	TemplateComponentSubtype string
	TemplateParameterType    string
	TemplateButtonPosition   *int
	InteractiveMessageType   string
	InteractiveHeaderType    string
	InteractiveButtonType    string
)

const (
	BlockingWait   BlockingType = "wait"
	BlockingNoWait BlockingType = "no_wait"
)

const (
	ContactStatusValid      ContactStatus = "valid"
	ContactStatusProcessing ContactStatus = "processing"
	ContactStatusInvalid    ContactStatus = "invalid"
	ContactStatusFailed     ContactStatus = "failed"
)

const (
	ContactAddressHome ContactAddressType = "HOME"
	ContactAddressWork ContactAddressType = "WORK"
)

const (
	ContactEmailHome ContactEmailType = "HOME"
	ContactEmailWork ContactEmailType = "WORK"
)

const (
	ContactPhoneCell   ContactPhoneType = "CELL"
	ContactPhoneMain   ContactPhoneType = "MAIN"
	ContactPhoneIphone ContactPhoneType = "IPHONE"
	ContactPhoneHome   ContactPhoneType = "HOME"
	ContactPhoneWork   ContactPhoneType = "WORK"
)

const (
	ContactURLHome ContactURLType = "HOME"
	ContactURLWork ContactURLType = "WORK"
)

type BaseResponse struct {
	Meta   *Metadata `json:"meta,omitempty"`
	Errors []Error   `json:"errors,omitempty"`
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
	Href    string `json:"href,omitempty"`
	Title   string `json:"title,omitempty"`
}

type Metadata struct {
	Success          bool   `json:"success,omitempty"`
	APIStatus        string `json:"api_status,omitempty"`
	Version          string `json:"version,omitempty"`
	HTTPCode         int    `json:"http_code,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
}

type ContactsRequest struct {
	Blocking   BlockingType `json:"blocking,omitempty" validate:"required,oneof=wait no_wait"`
	Contacts   []string     `json:"contacts,omitempty" validate:"required,min=1"`
	ForceCheck bool         `json:"force_check,omitempty"`
}

type ContactsResponse struct {
	BaseResponse
	Contacts []Contact `json:"contacts"`
}

type Contact struct {
	WaID   string        `json:"wa_id"`
	Input  string        `json:"input"`
	Status ContactStatus `json:"status"`
}

type Message struct {
	RecipientType RecipientType       `json:"recipient_type,omitempty"  validate:"required,eq=individual"`
	To            string              `json:"to" validate:"required,min=1"`
	Type          MessageType         `json:"type,omitempty" validate:"required,oneof=audio contact document image location sticker template text voice video interactive button"`
	Preview       bool                `json:"preview,omitempty"`
	Text          *MessageText        `json:"text,omitempty"`
	Audio         *MessageMedia       `json:"audio,omitempty"`
	Document      *MessageMedia       `json:"document,omitempty"`
	Image         *MessageMedia       `json:"image,omitempty"`
	Sticker       *MessageMedia       `json:"sticker,omitempty"`
	Video         *MessageMedia       `json:"video,omitempty"`
	Contacts      []MessageContact    `json:"contacts,omitempty"`
	Location      *MessageLocation    `json:"location,omitempty"`
	Template      *MessageTemplate    `json:"template,omitempty"`
	Interactive   *MessageInteractive `json:"interactive,omitempty"`
}

type InboundMessage struct {
	Message
	Button    *InboundMessageButton   `json:"button,omitempty"`
	Context   *InboundMessageContext  `json:"context,omitempty"`
	From      string                  `json:"from,omitempty"`
	ID        string                  `json:"id,omitempty"`
	Identity  *InboundMessageIdentity `json:"identity,omitempty"`
	Timestamp string                  `json:"timestamp,omitempty"`
	System    *MessageSystem          `json:"system,omitempty"`
	Voice     *MessageMedia           `json:"voice,omitempty"`
	Referral  *Referral               `json:"referral,omitempty"`
	Errors    []InboundError          `json:"errors,omitempty"`
}

type InboundWebhook struct {
	Contacts []InboundContact `json:"contacts,omitempty"`
	Messages []InboundMessage `json:"messages,omitempty"`
	Statuses []InboundStatus  `json:"statuses,omitempty"`
	Errors   []InboundError   `json:"errors,omitempty"`
}

type InboundContact struct {
	Profile *Profile `json:"profile,omitempty"`
	WaID    string   `json:"wa_id,omitempty"`
}

type InboundMessageButton struct {
	Payload string `json:"payload,omitempty"`
	Text    string `json:"text,omitempty"`
}

type InboundMessageIdentity struct {
	Acknowledged     string `json:"acknowledged,omitempty"`
	CreatedTimestamp string `json:"created_timestamp,omitempty"`
	Hash             string `json:"hash,omitempty"`
}

type InboundMessageContext struct {
	From                string   `json:"from,omitempty"`
	ID                  string   `json:"id,omitempty"`
	GroupID             string   `json:"group_id,omitempty"`
	Mentions            []string `json:"mentions,omitempty"`
	Forwarded           bool     `json:"forwarded,omitempty"`
	FrequentlyForwarded bool     `json:"frequently_forwarded,omitempty"`
}

type InboundStatus struct {
	Conversation *InboundStatusConversation `json:"conversation,omitempty"`
	ID           string                     `json:"id,omitempty"`
	Pricing      *InboundStatusPricing      `json:"pricing,omitempty"`
	RecipientID  string                     `json:"recipient_id,omitempty"`
	Status       string                     `json:"status,omitempty"`
	Timestamp    json.Number                `json:"timestamp,omitempty"`
	Type         string                     `json:"type,omitempty"`
}

type InboundStatusPricing struct {
	Billable     bool   `json:"billable,omitempty"`
	PricingModel string `json:"pricing_model,omitempty"`
}

type InboundStatusConversation struct {
	ID                  string                          `json:"id,omitempty"`
	Origin              InboundStatusConversationOrigin `json:"origin,omitempty"`
	ExpirationTimestamp json.Number                     `json:"expiration_timestamp,omitempty"`
}

type InboundStatusConversationOrigin struct {
	Type string `json:"type,omitempty"`
}

type InboundError struct {
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
	Title   string `json:"title,omitempty"`
}

type Referral struct {
	Headline   string        `json:"headline,omitempty"`
	Body       string        `json:"body,omitempty"`
	SourceType string        `json:"source_type,omitempty"`
	SourceID   string        `json:"source_id,omitempty"`
	SourceURL  string        `json:"source_url,omitempty"`
	Image      *MessageMedia `json:"image,omitempty"`
	Video      *MessageMedia `json:"video,omitempty"`
}

type MessageSystem struct {
	Body     string `json:"body,omitempty"`
	NewWaID  string `json:"new_wa_id,omitempty"`
	Type     string `json:"type,omitempty"`
	Identity string `json:"identity,omitempty"`
}

type Profile struct {
	Name string `json:"name,omitempty"`
}

type MessagesResponse struct {
	BaseResponse
	Messages []IDModel `json:"messages,omitempty"`
}

type IDModel struct {
	ID string `json:"id,omitempty"`
}

type MessageText struct {
	Body string `json:"body,omitempty"`
}

type MessageMedia struct {
	ID       string         `json:"id,omitempty"`
	Link     string         `json:"link,omitempty"`
	Caption  string         `json:"caption,omitempty"`
	Filename string         `json:"filename,omitempty"`
	Provider *MediaProvider `json:"provider,omitempty"`
	MIMEType string         `json:"mime_type,omitempty"`
	SHA256   string         `json:"sha256,omitempty"`
}

type MediaProvider struct {
	Name   string               `json:"name,omitempty"`
	Type   string               `json:"type,omitempty"`
	Config *MediaProviderConfig `json:"config,omitempty"`
}

type MediaProviderConfig struct {
	Bearer string                    `json:"bearer,omitempty"`
	Basic  *MediaProviderConfigBasic `json:"basic,omitempty"`
}

type MediaProviderConfigBasic struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type MessageContact struct {
	Addresses []ContactAddress `json:"addresses,omitempty"`
	Birthday  string           `json:"birthday,omitempty"`
	Emails    []ContactEmail   `json:"emails,omitempty"`
	Name      *ContactName     `json:"name,omitempty"`
	Org       *ContactOrg      `json:"org,omitempty"`
	Phones    []ContactPhone   `json:"phones,omitempty"`
	Urls      []ContactURL     `json:"urls,omitempty"`
}

type ContactAddress struct {
	Street      string             `json:"street,omitempty"`
	City        string             `json:"city,omitempty"`
	State       string             `json:"state,omitempty"`
	Zip         string             `json:"zip,omitempty"`
	Country     string             `json:"country,omitempty"`
	CountryCode string             `json:"country_code,omitempty"`
	Type        ContactAddressType `json:"type,omitempty"`
}

type ContactEmail struct {
	Email string           `json:"email,omitempty"`
	Type  ContactEmailType `json:"type,omitempty"`
}

type ContactName struct {
	FormattedName string `json:"formatted_name,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	MiddleName    string `json:"middle_name,omitempty"`
	Suffix        string `json:"suffix,omitempty"`
	Prefix        string `json:"prefix,omitempty"`
}

type ContactOrg struct {
	Company    string `json:"company,omitempty"`
	Department string `json:"department,omitempty"`
	Title      string `json:"title,omitempty"`
}

type ContactPhone struct {
	Phone string           `json:"phone,omitempty"`
	Type  ContactPhoneType `json:"type,omitempty"`
	WaID  string           `json:"wa_id,omitempty"`
}

type ContactURL struct {
	URL  string         `json:"url,omitempty"`
	Type ContactURLType `json:"type,omitempty"`
}

type MessageLocation struct {
	Longitude json.Number `json:"longitude,omitempty"`
	Latitude  json.Number `json:"latitude,omitempty"`
	Name      string      `json:"name,omitempty"`
	Address   string      `json:"address,omitempty"`
}

type MessageTemplate struct {
	Namespace  string              `json:"namespace,omitempty"`
	Name       string              `json:"name,omitempty"`
	Language   TemplateLanguage    `json:"language,omitempty"`
	Components []TemplateComponent `json:"components,omitempty"`
}

type TemplateLanguage struct {
	Policy TemplateLanguagePolicy `json:"policy,omitempty"`
	Code   string                 `json:"code,omitempty"`
}

type TemplateComponent struct {
	Type       TemplateComponentType    `json:"type,omitempty"`
	Subtype    TemplateComponentSubtype `json:"subtype,omitempty"`
	Parameters []TemplateParameter      `json:"parameters,omitempty"`
	Text       string                   `json:"text,omitempty"`
}

type TemplateParameter struct {
	Type    TemplateParameterType    `json:"type,omitempty"`
	SubType TemplateComponentSubtype `json:"sub_type,omitempty"`
	Index   TemplateButtonPosition   `json:"index,omitempty"`
	Caption string                   `json:"caption,omitempty"`
	Link    string                   `json:"link,omitempty"`
	Text    string                   `json:"text,omitempty"`
	Payload string                   `json:"payload,omitempty"`
}

type MessageInteractive struct {
	Type        InteractiveMessageType  `json:"type,omitempty"`
	Header      *InteractiveHeader      `json:"header,omitempty"`
	Body        *Text                   `json:"body,omitempty"`
	Footer      *Text                   `json:"footer,omitempty"`
	Action      *InteractiveAction      `json:"action,omitempty"`
	ListReply   *InteractiveSectionRow  `json:"list_reply,omitempty"`
	ButtonReply *InteractiveButtonReply `json:"button_reply,omitempty"`
}

type Text struct {
	Text string `json:"text,omitempty"`
}

type InteractiveHeader struct {
	Type     InteractiveHeaderType `json:"type,omitempty"`
	Text     string                `json:"text,omitempty"`
	Video    *MessageMedia         `json:"video,omitempty"`
	Image    *MessageMedia         `json:"image,omitempty"`
	Document *MessageMedia         `json:"document,omitempty"`
}

type InteractiveButton struct {
	Type  InteractiveButtonType   `json:"type,omitempty"`
	Reply *InteractiveButtonReply `json:"reply,omitempty"`
}

type InteractiveButtonReply struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type InteractiveSection struct {
	Title        string                      `json:"title,omitempty"`
	Rows         []InteractiveSectionRow     `json:"rows,omitempty"`
	ProductItems []InteractiveSectionProduct `json:"product_items,omitempty"`
}

type InteractiveSectionRow struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type InteractiveSectionProduct struct {
	ProductRetailerID string `json:"product_retailer_id,omitempty"`
}

type InteractiveAction struct {
	Button            string               `json:"button,omitempty"`
	Buttons           []InteractiveButton  `json:"buttons,omitempty"`
	CatalogID         string               `json:"catalog_id,omitempty"`
	Sections          []InteractiveSection `json:"sections,omitempty"`
	ProductRetailerID string               `json:"product_retailer_id,omitempty"`
}
