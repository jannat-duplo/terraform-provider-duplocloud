package duplocloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"terraform-provider-duplocloud/duplosdk"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	MAX_DUPLO_NO_HYPHEN_LENGTH          = len("duplo")
	MAX_DUPLO_LENGTH                    = len("duplo-")
	MAX_DUPLOSERVICES_LENGTH            = len("duploservices-1234567890ab-")
	MAX_DUPLOSERVICES_AND_SUFFIX_LENGTH = len("duploservices-1234567890ab--1234567890ab")
)

// Utility function to make a single schema element computed.
// Does not handle nested schema elements.
func makeSchemaComputed(el *schema.Schema) {
	el.Computed = true
	el.Required = false
	el.ForceNew = false
	el.Optional = false
	el.MaxItems = 0
	el.MinItems = 0
	el.Default = nil
	el.ValidateDiagFunc = nil
	el.ValidateFunc = nil
	el.DiffSuppressFunc = nil
	el.StateFunc = nil
	el.DefaultFunc = nil

	switch el.Elem.(type) {
	case *schema.Resource:
		for _, subel := range el.Elem.(*schema.Resource).Schema {
			makeSchemaComputed(subel)
		}
	default:
	}
}

// Utility function to convert the `from` interface to a JSON encoded string `field` in the `to` map.
func toJsonStringField(field string, from interface{}, to map[string]interface{}) {
	if json, err := json.Marshal(from); err == nil {
		to[field] = string(json)
	} else {
		log.Printf("[DEBUG] mapToJsonStringField: failed to serialize %s to JSON: %s", field, err)
	}
}

// Utility function to convert the `from` interface to a JSON encoded string `field`.
func toJsonStringState(field string, from interface{}, to *schema.ResourceData) {
	var err error
	var encoded []byte

	if encoded, err = json.Marshal(from); err == nil {
		err = to.Set(field, string(encoded))
	}

	if err != nil {
		log.Printf("[DEBUG] toJsonStringState: failed to serialize %s to JSON: %s", field, err)
	}
}

// Many kubernetes resources require a name to be a valid DNS subdomain, as defined in RFC 1123.
//nolint:staticcheck // TF needs to provide newer versions of these functions.
func ValidateDnsSubdomainRFC1123() schema.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(1, 253),
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9.-]*$`), "Invalid Kubernetes configmap name"),
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9]`), "Invalid Kubernetes configmap name"),
		validation.StringMatch(regexp.MustCompile(`[a-z0-9]$`), "Invalid Kubernetes configmap name"),
		validation.StringNotInSlice([]string{".."}, true),
	)
}

// ValidateJSONObjectString performs validation of a string that is supposed to be a JSON object.
func ValidateJSONArrayString(v interface{}, k string) (ws []string, errors []error) {
	// IAM Policy documents need to be valid JSON, and pass legacy parsing
	value := v.(string)
	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q contains invalid JSON", k))
		return
	}
	if value[:1] != "[" {
		errors = append(errors, fmt.Errorf("%q contains invalid JSON", k))
		return
	}
	if _, err := structure.NormalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

// ValidateJSONObjectString performs validation of a string that is supposed to be a JSON object.
func ValidateJSONObjectString(v interface{}, k string) (ws []string, errors []error) {
	// IAM Policy documents need to be valid JSON, and pass legacy parsing
	value := v.(string)
	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q contains invalid JSON", k))
		return
	}
	if value[:1] != "{" {
		errors = append(errors, fmt.Errorf("%q contains invalid JSON", k))
		return
	}
	if _, err := structure.NormalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

//nolint:deadcode,unused // utility function
func tagsSchemaForceNew() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func tagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

// awsTagsKeyValueSchema returns a Terraform schema to represent list of AWS tags.
//nolint:deadcode,unused // utility function
func awsTagsKeyValueSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

// awsTagsKeyValueSchema returns a Terraform schema to represent list of AWS tags.
//nolint:deadcode,unused // utility function
func awsTagsKeyValueSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 50,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringLenBetween(1, 128),
				},
				"value": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringLenBetween(0, 256),
				},
			},
		},
	}
}

// KeyValueSchema returns a Terraform schema to represent a key value pair
func KeyValueSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// CustomDataExSchema returns a Terraform schema to represent a key value pair with a type
func CustomDataExSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func keyValueToState(fieldName string, duploObjects *[]duplosdk.DuploKeyStringValue) []interface{} {
	if duploObjects != nil {
		input, _ := json.Marshal(&duploObjects)
		log.Printf("[TRACE] duplokeyValueToState[%s] ******** INPUT <= %s", fieldName, input)

		output := make([]interface{}, len(*duploObjects))
		for i, duploObject := range *duploObjects {
			jo := make(map[string]interface{})
			jo["key"] = duploObject.Key
			jo["value"] = duploObject.Value
			output[i] = jo
		}
		dump, _ := json.Marshal(output)
		log.Printf("[TRACE] duplokeyValueToState[%s] ******** OUTPUT => %s", fieldName, dump)
		return output
	}

	log.Printf("[TRACE] duplokeyValueToState[%s] ******** EMPTY INPUT", fieldName)
	return make([]interface{}, 0)
}

func keyValueFromState(fieldName string, d *schema.ResourceData) *[]duplosdk.DuploKeyStringValue {
	var ary []duplosdk.DuploKeyStringValue

	if v, ok := d.GetOk(fieldName); ok && v != nil && len(v.([]interface{})) > 0 {
		kvs := v.([]interface{})
		log.Printf("[TRACE] duploKeyValueFromState ********: found %s", fieldName)
		ary = make([]duplosdk.DuploKeyStringValue, 0, len(kvs))
		for _, raw := range kvs {
			kv := raw.(map[string]interface{})
			ary = append(ary, duplosdk.DuploKeyStringValue{
				Key:   kv["key"].(string),
				Value: kv["value"].(string),
			})
		}
	}

	return &ary
}

func customDataExToState(fieldName string, duploObjects *[]duplosdk.DuploCustomDataEx) []interface{} {
	if duploObjects != nil {
		input, _ := json.Marshal(&duploObjects)
		log.Printf("[TRACE] customDataExToState[%s] ******** INPUT <= %s", fieldName, input)

		output := make([]interface{}, len(*duploObjects))
		for i, duploObject := range *duploObjects {
			jo := make(map[string]interface{})
			jo["key"] = duploObject.Key
			jo["type"] = duploObject.Type
			jo["value"] = duploObject.Value
			output[i] = jo
		}
		dump, _ := json.Marshal(output)
		log.Printf("[TRACE] customDataExToState[%s] ******** OUTPUT => %s", fieldName, dump)
		return output
	}

	log.Printf("[TRACE] customDataExToState[%s] ******** EMPTY INPUT", fieldName)
	return make([]interface{}, 0)
}

func keyValueToMap(list *[]duplosdk.DuploKeyStringValue) map[string]interface{} {
	result := map[string]interface{}{}
	if list != nil {
		for _, item := range *list {
			result[item.Key] = item.Value
		}
	}
	return result
}

func keyValueFromMap(d map[string]interface{}) *[]duplosdk.DuploKeyStringValue {
	list := make([]duplosdk.DuploKeyStringValue, 0, len(d))

	for k, v := range d {
		list = append(list, duplosdk.DuploKeyStringValue{
			Key:   k,
			Value: v.(string),
		})
	}

	return &list
}

func keyValueFromStateList(fieldName string, d map[string]interface{}) *[]duplosdk.DuploKeyStringValue {
	var ary []duplosdk.DuploKeyStringValue

	if v, ok := d[fieldName]; ok && v != nil && len(v.([]interface{})) > 0 {
		kvs := v.([]interface{})
		log.Printf("[TRACE] duploKeyValueFromMap ********: found %s", fieldName)
		ary = make([]duplosdk.DuploKeyStringValue, 0, len(kvs))
		for _, raw := range kvs {
			kv := raw.(map[string]interface{})
			ary = append(ary, duplosdk.DuploKeyStringValue{
				Key:   kv["key"].(string),
				Value: kv["value"].(string),
			})
		}
	}

	return &ary
}

// FilterSchema returns a Terraform schema to represent a filter
func FilterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
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
	}
}

// FiltersSchema returns a Terraform schema to represent a set of filters
func FiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"values": {
					Type:     schema.TypeString,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

// Utility function to return a pointer to a single valid (but optional) resource data block.
func getOptionalBlock(data *schema.ResourceData, key string) (*interface{}, error) {
	var value *interface{}

	if v, ok := data.GetOk(key); ok {
		x := v.([]interface{})

		if len(x) == 1 {
			if x[0] == nil {
				return nil, fmt.Errorf("at least one field is expected inside %s", key)
			}
			value = &x[0]
		}
	}

	return value, nil
}

// Utility function to return a pointer to a single valid (but optional) resource data block as a map.
func getOptionalBlockAsMap(data *schema.ResourceData, key string) (map[string]interface{}, error) {
	block, err := getOptionalBlock(data, key)
	if block == nil || err != nil {
		return nil, err
	}
	return (*block).(map[string]interface{}), nil
}

// Utility function to return a pointer to a single valid (but optional) resource data block.
func getOptionalNestedBlock(d map[string]interface{}, key string) (*interface{}, error) {
	var value *interface{}

	if v, ok := d[key]; ok {
		x := v.([]interface{})

		if len(x) == 1 {
			if x[0] == nil {
				return nil, fmt.Errorf("at least one field is expected inside %s", key)
			}
			value = &x[0]
		}
	}

	return value, nil
}

// Utility function to return a pointer to a single valid (but optional) resource data block as a map.
func getOptionalNestedBlockAsMap(d map[string]interface{}, key string) (map[string]interface{}, error) {
	block, err := getOptionalNestedBlock(d, key)
	if block == nil || err != nil {
		return nil, err
	}
	return (*block).(map[string]interface{}), nil
}

func getAsStringArray(data *schema.ResourceData, key string) (*[]string, bool) {
	var ok bool
	var result []string
	var v interface{}

	if v, ok = data.GetOk(key); ok && v != nil {
		list := v.([]interface{})
		result = make([]string, len(list))
		for i, el := range list {
			result[i] = el.(string)
		}
	}

	return &result, ok
}

// Utiliy function to return a filtered list of tenant metadata, given the selected keys.
func selectKeyValuesFromMap(metadata *[]duplosdk.DuploKeyStringValue, keys map[string]interface{}) *[]duplosdk.DuploKeyStringValue {
	settings := make([]duplosdk.DuploKeyStringValue, 0, len(keys))
	for _, kv := range *metadata {
		if _, ok := keys[kv.Key]; ok {
			settings = append(settings, kv)
		}
	}

	return &settings
}

// Utiliy function to return a filtered list of tenant metadata, given the selected keys.
func selectKeyValues(metadata *[]duplosdk.DuploKeyStringValue, keys []string) *[]duplosdk.DuploKeyStringValue {
	specified := map[string]interface{}{}
	for _, k := range keys {
		specified[k] = struct{}{}
	}

	return selectKeyValuesFromMap(metadata, specified)
}

// Internal function used to re-order key value pairs
//nolint:deadcode,unused // utility function
func reorderKeyValues(pairs []interface{}) {

	// Re-order environment variables to a canonical order.
	sort.SliceStable(pairs, func(i, j int) bool {
		mi := pairs[i].(map[string]interface{})
		mj := pairs[j].(map[string]interface{})

		// Get both name keys, fall back on an empty string.
		si := ""
		sj := ""
		if v, ok := mi["key"]; ok && !isInterfaceNil(v) {
			si = v.(string)
		}
		if v, ok := mj["key"]; ok && !isInterfaceNil(v) {
			sj = v.(string)
		}

		// Compare the two.
		return si < sj
	})
}

func getStringArray(data map[string]interface{}, key string) (*[]string, bool) {
	var ok bool
	var result []string
	var v interface{}

	if v, ok = data[key]; ok && v != nil {
		list := v.([]interface{})
		result = make([]string, len(list))
		for i, el := range list {
			result[i] = el.(string)
		}
	}

	return &result, ok
}

func waitForResourceToBeMissingAfterDelete(ctx context.Context, d *schema.ResourceData, kind string, id string, get func() (interface{}, duplosdk.ClientError)) diag.Diagnostics {
	err := resource.RetryContext(ctx, d.Timeout("delete"), func() *resource.RetryError {
		resp, errget := get()

		if errget != nil {
			if errget.Status() == 404 {
				return nil
			}

			return resource.NonRetryableError(fmt.Errorf("error getting %s '%s': %s", kind, id, errget))
		}

		if !isInterfaceNil(resp) {
			return resource.RetryableError(fmt.Errorf("expected %s '%s' to be missing, but it still exists", kind, id))
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("error deleting %s '%s': %s", kind, id, err)
	}
	return nil
}

func waitForResourceToBePresentAfterCreate(ctx context.Context, d *schema.ResourceData, kind string, id string, get func() (interface{}, duplosdk.ClientError)) diag.Diagnostics {
	err := resource.RetryContext(ctx, d.Timeout("create"), func() *resource.RetryError {
		resp, errget := get()

		if errget != nil {
			if errget.Status() == 404 {
				return resource.RetryableError(fmt.Errorf("expected %s '%s' to be retrieved, but got a 404", kind, id))
			}

			return resource.NonRetryableError(fmt.Errorf("error getting %s '%s': %s", kind, id, errget))
		}

		if isInterfaceNil(resp) {
			return resource.RetryableError(fmt.Errorf("expected %s '%s' to be retrieved, but got: nil", kind, id))
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("error creating %s '%s': %s", kind, id, err)
	}
	return nil
}

func isInterfaceNil(v interface{}) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

func isInterfaceNilOrEmptySlice(v interface{}) bool {
	if isInterfaceNil(v) {
		return true
	}

	slice := reflect.ValueOf(v)
	return slice.Kind() == reflect.Slice && slice.IsValid() && (slice.IsNil() || slice.Len() == 0)
}

func isInterfaceEmptySlice(v interface{}) bool {
	slice := reflect.ValueOf(v)
	return slice.Kind() == reflect.Slice && slice.IsValid() && !slice.IsNil() && slice.Len() == 0
}

//nolint:deadcode,unused // utility function
func isInterfaceEmptyMap(v interface{}) bool {
	emap := reflect.ValueOf(v)
	return emap.Kind() == reflect.Map && emap.IsValid() && !emap.IsNil() && emap.Len() == 0
}

func isInterfaceNilOrEmptyMap(v interface{}) bool {
	if isInterfaceNil(v) {
		return true
	}

	emap := reflect.ValueOf(v)
	return emap.Kind() == reflect.Map && emap.IsValid() && (emap.IsNil() || emap.Len() == 0)
}

// Internal function to check if a given encoded JSON value represents a valid JSON object array.
func validateJsonObjectArray(key string, value string) (ws []string, errors []error) {
	result := []map[string]interface{}{}
	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		errors = append(errors, fmt.Errorf("%s is invalid: %s", key, err))
	}
	return
}

// Internal function to convert map keys from lower camel-case to upper camel-case.
//  - Adds an upper camel-case entry for each lower camel-case entry, unless the upper exists already.
//  - Removes any lower camel-case entry.
//  - Never overwrites any existing upper camel-case keys.
func makeMapUpperCamelCase(m map[string]interface{}) {
	for k := range m {

		// Only convert lowercase entries.
		if unicode.IsLower([]rune(k)[0]) {
			//nolint:staticcheck // SA1019 ignore this!
			upper := strings.Title(k)

			// Add the upper camel-case entry, if it doesn't exist.
			if _, ok := m[upper]; !ok {
				m[upper] = m[k]
			}

			// Remove the lower camel-case entry.
			delete(m, k)
		}
	}
}

// Internal function to reduce empty or nil map entries.
func reduceNilOrEmptyMapEntries(m map[string]interface{}) {
	for k, v := range m {
		if isInterfaceNil(v) || isInterfaceNilOrEmptyMap(v) || isInterfaceNilOrEmptySlice(v) {
			delete(m, k)
		}
	}
}

func waitForResourceWithStatusDone(ctx context.Context, d *schema.ResourceData, kind string, id string, get func() (bool, duplosdk.ClientError)) diag.Diagnostics {
	err := resource.RetryContext(ctx, d.Timeout(kind), func() *resource.RetryError {
		status, errget := get()

		if errget != nil {
			if errget.Status() == 404 {
				return nil
			}

			return resource.NonRetryableError(fmt.Errorf("error getting %s '%s': %s", kind, id, errget))
		}
		// return nil if want to complete wait
		if !status {
			return resource.RetryableError(fmt.Errorf("expected %s '%s' to be missing, but it still exists", kind, id))
		}

		return nil
	})
	if err != nil {
		return diag.Errorf("error deleting %s '%s': %s", kind, id, err)
	}
	return nil
}

func flattenStringMap(duplo map[string]string) map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range duplo {
		m[k] = v
	}
	return m
}

func flattenGcpLabels(d *schema.ResourceData, duplo map[string]string) {
	d.Set("labels", flattenStringMap(duplo))
}

func expandAsStringMap(fieldName string, d *schema.ResourceData) map[string]string {
	m := map[string]string{}

	if v, ok := d.GetOk(fieldName); ok && v != nil && len(v.(map[string]interface{})) > 0 {
		for k, v := range v.(map[string]interface{}) {
			if v == nil {
				m[k] = ""
			} else {
				m[k] = v.(string)
			}
		}
	}

	return m
}

func expandStringMap(fieldName string, d map[string]interface{}) map[string]string {
	m := map[string]string{}

	if v, ok := d[fieldName]; ok && v != nil && len(v.(map[string]interface{})) > 0 {
		for k, v := range v.(map[string]interface{}) {
			if v == nil {
				m[k] = ""
			} else {
				m[k] = v.(string)
			}
		}
	}

	return m
}

func errorsToDiagnostics(prefix string, errors []error) diag.Diagnostics {
	if len(errors) == 0 {
		return nil
	}

	diags := make(diag.Diagnostics, 0, len(errors))
	for i := range errors {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s%s", prefix, errors[i]),
		})
	}
	return diags
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func flattenStringSet(list []string) *schema.Set {
	return schema.NewSet(schema.HashString, flattenStringList(list))
}

func objectMapToStringMap(rm map[string]interface{}) map[string]string {
	result := map[string]string{}
	for k, v := range rm {
		switch assertedValue := v.(type) {
		case string:
			result[k] = assertedValue
		default:
			// Make a best effort to coerce into a string, even if underlying type is not a string
			log.Printf("[DEBUG] non-string value encountered for key '%s' while converting object map to string map", k)
			result[k] = fmt.Sprintf("%v", assertedValue)
		}
	}
	return result
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

func expandStringSet(configured *schema.Set) []string {
	return expandStringList(configured.List())
}

func CaseDifference(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

func HashStringIgnoreCase(v interface{}) int {
	return schema.HashString(strings.ToLower(v.(string)))
}

func Base64EncodeIfNot(data string) string {
	// Check whether the data is already Base64 encoded; don't double-encode
	if base64IsEncoded(data) {
		return data
	}
	// data has not been encoded encode and return
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func base64IsEncoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}

func DuploManagedAzureTags() []string {
	return []string{"TENANT_NAME", "TENANT_ID", "duplo-project", "duplo_creation_time", "duplo_sync_vm", "owner"}
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
