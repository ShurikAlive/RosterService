# {{classname}}

All URIs are relative to *http://localhost:8181/*

Method | HTTP request | Description
------------- | ------------- | -------------
[**EquipmentEquipmentIdDelete**](WargearsApi.md#EquipmentEquipmentIdDelete) | **Delete** /equipment/{equipmentId} | Delete equipment by ID.
[**EquipmentEquipmentIdGet**](WargearsApi.md#EquipmentEquipmentIdGet) | **Get** /equipment/{equipmentId} | Returns a equipment by ID.
[**EquipmentEquipmentIdPut**](WargearsApi.md#EquipmentEquipmentIdPut) | **Put** /equipment/{equipmentId} | Edit equipment by ID.
[**EquipmentGet**](WargearsApi.md#EquipmentGet) | **Get** /equipment | Returns a list of available equipment.
[**EquipmentPost**](WargearsApi.md#EquipmentPost) | **Post** /equipment | Add new equipment.

# **EquipmentEquipmentIdDelete**
> string EquipmentEquipmentIdDelete(ctx, equipmentId)
Delete equipment by ID.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **equipmentId** | **string**| ID Equipment | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **EquipmentEquipmentIdGet**
> Equipment EquipmentEquipmentIdGet(ctx, equipmentId)
Returns a equipment by ID.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **equipmentId** | **string**| ID Equipment | 

### Return type

[**Equipment**](Equipment.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **EquipmentEquipmentIdPut**
> string EquipmentEquipmentIdPut(ctx, body, equipmentId)
Edit equipment by ID.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**EditEquipment**](EditEquipment.md)|  | 
  **equipmentId** | **string**| ID Equipment | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **EquipmentGet**
> []Equipment EquipmentGet(ctx, )
Returns a list of available equipment.

Return JSON array of info available equipments.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]Equipment**](Equipment.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **EquipmentPost**
> string EquipmentPost(ctx, body)
Add new equipment.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**EditEquipment**](EditEquipment.md)|  | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

