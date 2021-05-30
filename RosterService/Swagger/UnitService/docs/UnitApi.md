# {{classname}}

All URIs are relative to *http://localhost:8181/*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UnitGet**](UnitApi.md#UnitGet) | **Get** /unit | Returns a list of available units.
[**UnitPost**](UnitApi.md#UnitPost) | **Post** /unit | Add new unit.
[**UnitUnitIdDelete**](UnitApi.md#UnitUnitIdDelete) | **Delete** /unit/{unitId} | Delete unit by ID.
[**UnitUnitIdGet**](UnitApi.md#UnitUnitIdGet) | **Get** /unit/{unitId} | Returns a unit by ID.
[**UnitUnitIdPut**](UnitApi.md#UnitUnitIdPut) | **Put** /unit/{unitId} | Edit unit by ID.

# **UnitGet**
> []Unit UnitGet(ctx, )
Returns a list of available units.

Return JSON array of info available units.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]Unit**](Unit.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UnitPost**
> string UnitPost(ctx, body)
Add new unit.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**EditUnit**](EditUnit.md)|  | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UnitUnitIdDelete**
> string UnitUnitIdDelete(ctx, unitId)
Delete unit by ID.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **unitId** | **string**| ID unit | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UnitUnitIdGet**
> Unit UnitUnitIdGet(ctx, unitId)
Returns a unit by ID.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **unitId** | **string**| ID unit | 

### Return type

[**Unit**](Unit.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UnitUnitIdPut**
> string UnitUnitIdPut(ctx, body, unitId)
Edit unit by ID.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**EditUnit**](EditUnit.md)|  | 
  **unitId** | **string**| ID unit | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

