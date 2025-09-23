# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/v1/common.proto](#api_v1_common-proto)
    - [State](#monotreme-api-v1-State)
    - [Visibility](#monotreme-api-v1-Visibility)
  
- [api/v1/user_service.proto](#api_v1_user_service-proto)
    - [CreateUserAccessTokenRequest](#monotreme-api-v1-CreateUserAccessTokenRequest)
    - [CreateUserRequest](#monotreme-api-v1-CreateUserRequest)
    - [DeleteUserAccessTokenRequest](#monotreme-api-v1-DeleteUserAccessTokenRequest)
    - [DeleteUserRequest](#monotreme-api-v1-DeleteUserRequest)
    - [GetUserRequest](#monotreme-api-v1-GetUserRequest)
    - [ListUserAccessTokensRequest](#monotreme-api-v1-ListUserAccessTokensRequest)
    - [ListUserAccessTokensResponse](#monotreme-api-v1-ListUserAccessTokensResponse)
    - [ListUsersRequest](#monotreme-api-v1-ListUsersRequest)
    - [ListUsersResponse](#monotreme-api-v1-ListUsersResponse)
    - [UpdateUserRequest](#monotreme-api-v1-UpdateUserRequest)
    - [User](#monotreme-api-v1-User)
    - [UserAccessToken](#monotreme-api-v1-UserAccessToken)
  
    - [Role](#monotreme-api-v1-Role)
  
    - [UserService](#monotreme-api-v1-UserService)
  
- [api/v1/auth_service.proto](#api_v1_auth_service-proto)
    - [GetAuthStatusRequest](#monotreme-api-v1-GetAuthStatusRequest)
    - [SignInRequest](#monotreme-api-v1-SignInRequest)
    - [SignInWithSSORequest](#monotreme-api-v1-SignInWithSSORequest)
    - [SignOutRequest](#monotreme-api-v1-SignOutRequest)
    - [SignUpRequest](#monotreme-api-v1-SignUpRequest)
  
    - [AuthService](#monotreme-api-v1-AuthService)
  
- [api/v1/collection_service.proto](#api_v1_collection_service-proto)
    - [Collection](#monotreme-api-v1-Collection)
    - [CreateCollectionRequest](#monotreme-api-v1-CreateCollectionRequest)
    - [DeleteCollectionRequest](#monotreme-api-v1-DeleteCollectionRequest)
    - [GetCollectionByNameRequest](#monotreme-api-v1-GetCollectionByNameRequest)
    - [GetCollectionRequest](#monotreme-api-v1-GetCollectionRequest)
    - [ImportBookmarksRequest](#monotreme-api-v1-ImportBookmarksRequest)
    - [ImportBookmarksResponse](#monotreme-api-v1-ImportBookmarksResponse)
    - [ListCollectionsRequest](#monotreme-api-v1-ListCollectionsRequest)
    - [ListCollectionsResponse](#monotreme-api-v1-ListCollectionsResponse)
    - [UpdateCollectionRequest](#monotreme-api-v1-UpdateCollectionRequest)
  
    - [CollectionService](#monotreme-api-v1-CollectionService)
  
- [api/v1/shortcut_service.proto](#api_v1_shortcut_service-proto)
    - [CreateShortcutRequest](#monotreme-api-v1-CreateShortcutRequest)
    - [DeleteShortcutRequest](#monotreme-api-v1-DeleteShortcutRequest)
    - [GetShortcutAnalyticsRequest](#monotreme-api-v1-GetShortcutAnalyticsRequest)
    - [GetShortcutAnalyticsResponse](#monotreme-api-v1-GetShortcutAnalyticsResponse)
    - [GetShortcutAnalyticsResponse.AnalyticsItem](#monotreme-api-v1-GetShortcutAnalyticsResponse-AnalyticsItem)
    - [GetShortcutByNameRequest](#monotreme-api-v1-GetShortcutByNameRequest)
    - [GetShortcutRequest](#monotreme-api-v1-GetShortcutRequest)
    - [ListShortcutsRequest](#monotreme-api-v1-ListShortcutsRequest)
    - [ListShortcutsResponse](#monotreme-api-v1-ListShortcutsResponse)
    - [Shortcut](#monotreme-api-v1-Shortcut)
    - [Shortcut.OpenGraphMetadata](#monotreme-api-v1-Shortcut-OpenGraphMetadata)
    - [UpdateShortcutRequest](#monotreme-api-v1-UpdateShortcutRequest)
  
    - [ShortcutService](#monotreme-api-v1-ShortcutService)
  
- [api/v1/subscription_service.proto](#api_v1_subscription_service-proto)
    - [DeleteSubscriptionRequest](#monotreme-api-v1-DeleteSubscriptionRequest)
    - [GetSubscriptionRequest](#monotreme-api-v1-GetSubscriptionRequest)
    - [Subscription](#monotreme-api-v1-Subscription)
    - [UpdateSubscriptionRequest](#monotreme-api-v1-UpdateSubscriptionRequest)
  
    - [PlanType](#monotreme-api-v1-PlanType)
  
    - [SubscriptionService](#monotreme-api-v1-SubscriptionService)
  
- [api/v1/user_setting_service.proto](#api_v1_user_setting_service-proto)
    - [GetUserSettingRequest](#monotreme-api-v1-GetUserSettingRequest)
    - [UpdateUserSettingRequest](#monotreme-api-v1-UpdateUserSettingRequest)
    - [UserSetting](#monotreme-api-v1-UserSetting)
    - [UserSetting.AccessTokensSetting](#monotreme-api-v1-UserSetting-AccessTokensSetting)
    - [UserSetting.AccessTokensSetting.AccessToken](#monotreme-api-v1-UserSetting-AccessTokensSetting-AccessToken)
    - [UserSetting.GeneralSetting](#monotreme-api-v1-UserSetting-GeneralSetting)
  
    - [UserSettingService](#monotreme-api-v1-UserSettingService)
  
- [api/v1/workspace_service.proto](#api_v1_workspace_service-proto)
    - [GetWorkspaceProfileRequest](#monotreme-api-v1-GetWorkspaceProfileRequest)
    - [GetWorkspaceSettingRequest](#monotreme-api-v1-GetWorkspaceSettingRequest)
    - [IdentityProvider](#monotreme-api-v1-IdentityProvider)
    - [IdentityProviderConfig](#monotreme-api-v1-IdentityProviderConfig)
    - [IdentityProviderConfig.FieldMapping](#monotreme-api-v1-IdentityProviderConfig-FieldMapping)
    - [IdentityProviderConfig.OAuth2Config](#monotreme-api-v1-IdentityProviderConfig-OAuth2Config)
    - [UpdateWorkspaceSettingRequest](#monotreme-api-v1-UpdateWorkspaceSettingRequest)
    - [WorkspaceProfile](#monotreme-api-v1-WorkspaceProfile)
    - [WorkspaceSetting](#monotreme-api-v1-WorkspaceSetting)
  
    - [IdentityProvider.Type](#monotreme-api-v1-IdentityProvider-Type)
  
    - [WorkspaceService](#monotreme-api-v1-WorkspaceService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_v1_common-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v1/common.proto


 


<a name="monotreme-api-v1-State"></a>

### State


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_UNSPECIFIED | 0 |  |
| ACTIVE | 1 |  |
| INACTIVE | 2 |  |



<a name="monotreme-api-v1-Visibility"></a>

### Visibility


| Name | Number | Description |
| ---- | ------ | ----------- |
| VISIBILITY_UNSPECIFIED | 0 |  |
| WORKSPACE | 1 |  |
| PUBLIC | 2 |  |


 

 

 



<a name="api_v1_user_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v1/user_service.proto



<a name="monotreme-api-v1-CreateUserAccessTokenRequest"></a>

### CreateUserAccessTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  | id is the user id. |
| description | [string](#string) |  | description is the description of the access token. |
| expires_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional | expires_at is the expiration time of the access token. If expires_at is not set, the access token will never expire. |






<a name="monotreme-api-v1-CreateUserRequest"></a>

### CreateUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#monotreme-api-v1-User) |  |  |






<a name="monotreme-api-v1-DeleteUserAccessTokenRequest"></a>

### DeleteUserAccessTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  | id is the user id. |
| access_token | [string](#string) |  | access_token is the access token to delete. |






<a name="monotreme-api-v1-DeleteUserRequest"></a>

### DeleteUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |






<a name="monotreme-api-v1-GetUserRequest"></a>

### GetUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |






<a name="monotreme-api-v1-ListUserAccessTokensRequest"></a>

### ListUserAccessTokensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  | id is the user id. |






<a name="monotreme-api-v1-ListUserAccessTokensResponse"></a>

### ListUserAccessTokensResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_tokens | [UserAccessToken](#monotreme-api-v1-UserAccessToken) | repeated |  |






<a name="monotreme-api-v1-ListUsersRequest"></a>

### ListUsersRequest







<a name="monotreme-api-v1-ListUsersResponse"></a>

### ListUsersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users | [User](#monotreme-api-v1-User) | repeated |  |






<a name="monotreme-api-v1-UpdateUserRequest"></a>

### UpdateUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#monotreme-api-v1-User) |  |  |
| update_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  |  |






<a name="monotreme-api-v1-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |
| state | [State](#monotreme-api-v1-State) |  |  |
| created_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| updated_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| role | [Role](#monotreme-api-v1-Role) |  |  |
| email | [string](#string) |  |  |
| nickname | [string](#string) |  |  |
| password | [string](#string) |  |  |






<a name="monotreme-api-v1-UserAccessToken"></a>

### UserAccessToken



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| description | [string](#string) |  |  |
| issued_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| expires_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |





 


<a name="monotreme-api-v1-Role"></a>

### Role


| Name | Number | Description |
| ---- | ------ | ----------- |
| ROLE_UNSPECIFIED | 0 |  |
| ADMIN | 1 |  |
| USER | 2 |  |


 

 


<a name="monotreme-api-v1-UserService"></a>

### UserService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListUsers | [ListUsersRequest](#monotreme-api-v1-ListUsersRequest) | [ListUsersResponse](#monotreme-api-v1-ListUsersResponse) | ListUsers returns a list of users. |
| GetUser | [GetUserRequest](#monotreme-api-v1-GetUserRequest) | [User](#monotreme-api-v1-User) | GetUser returns a user by id. |
| CreateUser | [CreateUserRequest](#monotreme-api-v1-CreateUserRequest) | [User](#monotreme-api-v1-User) | CreateUser creates a new user. |
| UpdateUser | [UpdateUserRequest](#monotreme-api-v1-UpdateUserRequest) | [User](#monotreme-api-v1-User) |  |
| DeleteUser | [DeleteUserRequest](#monotreme-api-v1-DeleteUserRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) | DeleteUser deletes a user by id. |
| ListUserAccessTokens | [ListUserAccessTokensRequest](#monotreme-api-v1-ListUserAccessTokensRequest) | [ListUserAccessTokensResponse](#monotreme-api-v1-ListUserAccessTokensResponse) | ListUserAccessTokens returns a list of access tokens for a user. |
| CreateUserAccessToken | [CreateUserAccessTokenRequest](#monotreme-api-v1-CreateUserAccessTokenRequest) | [UserAccessToken](#monotreme-api-v1-UserAccessToken) | CreateUserAccessToken creates a new access token for a user. |
| DeleteUserAccessToken | [DeleteUserAccessTokenRequest](#monotreme-api-v1-DeleteUserAccessTokenRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) | DeleteUserAccessToken deletes an access token for a user. |

 



<a name="api_v1_auth_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v1/auth_service.proto



<a name="monotreme-api-v1-GetAuthStatusRequest"></a>

### GetAuthStatusRequest







<a name="monotreme-api-v1-SignInRequest"></a>

### SignInRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [string](#string) |  |  |
| password | [string](#string) |  |  |






<a name="monotreme-api-v1-SignInWithSSORequest"></a>

### SignInWithSSORequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| idp_id | [string](#string) |  | The id of the SSO provider. |
| code | [string](#string) |  | The code to sign in with. |
| redirect_uri | [string](#string) |  | The redirect URI. |






<a name="monotreme-api-v1-SignOutRequest"></a>

### SignOutRequest







<a name="monotreme-api-v1-SignUpRequest"></a>

### SignUpRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [string](#string) |  |  |
| nickname | [string](#string) |  |  |
| password | [string](#string) |  |  |





 

 

 


<a name="monotreme-api-v1-AuthService"></a>

### AuthService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAuthStatus | [GetAuthStatusRequest](#monotreme-api-v1-GetAuthStatusRequest) | [User](#monotreme-api-v1-User) | GetAuthStatus returns the current auth status of the user. |
| SignIn | [SignInRequest](#monotreme-api-v1-SignInRequest) | [User](#monotreme-api-v1-User) | SignIn signs in the user with the given username and password. |
| SignInWithSSO | [SignInWithSSORequest](#monotreme-api-v1-SignInWithSSORequest) | [User](#monotreme-api-v1-User) | SignInWithSSO signs in the user with the given SSO code. |
| SignUp | [SignUpRequest](#monotreme-api-v1-SignUpRequest) | [User](#monotreme-api-v1-User) | SignUp signs up the user with the given username and password. |
| SignOut | [SignOutRequest](#monotreme-api-v1-SignOutRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) | SignOut signs out the user. |

 



<a name="api_v1_collection_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v1/collection_service.proto



<a name="monotreme-api-v1-Collection"></a>

### Collection



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |
| creator_id | [int32](#int32) |  |  |
| created_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| updated_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| name | [string](#string) |  |  |
| title | [string](#string) |  |  |
| description | [string](#string) |  |  |
| shortcut_ids | [int32](#int32) | repeated |  |
| visibility | [Visibility](#monotreme-api-v1-Visibility) |  |  |






<a name="monotreme-api-v1-CreateCollectionRequest"></a>

### CreateCollectionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| collection | [Collection](#monotreme-api-v1-Collection) |  |  |






<a name="monotreme-api-v1-DeleteCollectionRequest"></a>

### DeleteCollectionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |






<a name="monotreme-api-v1-GetCollectionByNameRequest"></a>

### GetCollectionByNameRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="monotreme-api-v1-GetCollectionRequest"></a>

### GetCollectionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |






<a name="monotreme-api-v1-ImportBookmarksRequest"></a>

### ImportBookmarksRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| html_content | [string](#string) |  |  |






<a name="monotreme-api-v1-ImportBookmarksResponse"></a>

### ImportBookmarksResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| collections | [Collection](#monotreme-api-v1-Collection) | repeated |  |
| total_shortcuts | [int32](#int32) |  |  |
| total_collections | [int32](#int32) |  |  |






<a name="monotreme-api-v1-ListCollectionsRequest"></a>

### ListCollectionsRequest







<a name="monotreme-api-v1-ListCollectionsResponse"></a>

### ListCollectionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| collections | [Collection](#monotreme-api-v1-Collection) | repeated |  |






<a name="monotreme-api-v1-UpdateCollectionRequest"></a>

### UpdateCollectionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| collection | [Collection](#monotreme-api-v1-Collection) |  |  |
| update_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  |  |





 

 

 


<a name="monotreme-api-v1-CollectionService"></a>

### CollectionService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListCollections | [ListCollectionsRequest](#monotreme-api-v1-ListCollectionsRequest) | [ListCollectionsResponse](#monotreme-api-v1-ListCollectionsResponse) | ListCollections returns a list of collections. |
| GetCollection | [GetCollectionRequest](#monotreme-api-v1-GetCollectionRequest) | [Collection](#monotreme-api-v1-Collection) | GetCollection returns a collection by id. |
| GetCollectionByName | [GetCollectionByNameRequest](#monotreme-api-v1-GetCollectionByNameRequest) | [Collection](#monotreme-api-v1-Collection) | GetCollectionByName returns a collection by name. |
| CreateCollection | [CreateCollectionRequest](#monotreme-api-v1-CreateCollectionRequest) | [Collection](#monotreme-api-v1-Collection) | CreateCollection creates a collection. |
| UpdateCollection | [UpdateCollectionRequest](#monotreme-api-v1-UpdateCollectionRequest) | [Collection](#monotreme-api-v1-Collection) | UpdateCollection updates a collection. |
| DeleteCollection | [DeleteCollectionRequest](#monotreme-api-v1-DeleteCollectionRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) | DeleteCollection deletes a collection by id. |
| ImportBookmarks | [ImportBookmarksRequest](#monotreme-api-v1-ImportBookmarksRequest) | [ImportBookmarksResponse](#monotreme-api-v1-ImportBookmarksResponse) | ImportBookmarks imports bookmarks from an HTML file and creates collections and shortcuts. |

 



<a name="api_v1_shortcut_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v1/shortcut_service.proto



<a name="monotreme-api-v1-CreateShortcutRequest"></a>

### CreateShortcutRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| shortcut | [Shortcut](#monotreme-api-v1-Shortcut) |  |  |






<a name="monotreme-api-v1-DeleteShortcutRequest"></a>

### DeleteShortcutRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |






<a name="monotreme-api-v1-GetShortcutAnalyticsRequest"></a>

### GetShortcutAnalyticsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |






<a name="monotreme-api-v1-GetShortcutAnalyticsResponse"></a>

### GetShortcutAnalyticsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| references | [GetShortcutAnalyticsResponse.AnalyticsItem](#monotreme-api-v1-GetShortcutAnalyticsResponse-AnalyticsItem) | repeated |  |
| devices | [GetShortcutAnalyticsResponse.AnalyticsItem](#monotreme-api-v1-GetShortcutAnalyticsResponse-AnalyticsItem) | repeated |  |
| browsers | [GetShortcutAnalyticsResponse.AnalyticsItem](#monotreme-api-v1-GetShortcutAnalyticsResponse-AnalyticsItem) | repeated |  |






<a name="monotreme-api-v1-GetShortcutAnalyticsResponse-AnalyticsItem"></a>

### GetShortcutAnalyticsResponse.AnalyticsItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| count | [int32](#int32) |  |  |






<a name="monotreme-api-v1-GetShortcutByNameRequest"></a>

### GetShortcutByNameRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |






<a name="monotreme-api-v1-GetShortcutRequest"></a>

### GetShortcutRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |






<a name="monotreme-api-v1-ListShortcutsRequest"></a>

### ListShortcutsRequest







<a name="monotreme-api-v1-ListShortcutsResponse"></a>

### ListShortcutsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| shortcuts | [Shortcut](#monotreme-api-v1-Shortcut) | repeated |  |






<a name="monotreme-api-v1-Shortcut"></a>

### Shortcut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |
| uuid | [string](#string) |  |  |
| creator_id | [int32](#int32) |  |  |
| created_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| updated_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| name | [string](#string) |  |  |
| link | [string](#string) |  |  |
| title | [string](#string) |  |  |
| tags | [string](#string) | repeated |  |
| description | [string](#string) |  |  |
| visibility | [Visibility](#monotreme-api-v1-Visibility) |  |  |
| view_count | [int32](#int32) |  |  |
| og_metadata | [Shortcut.OpenGraphMetadata](#monotreme-api-v1-Shortcut-OpenGraphMetadata) |  |  |






<a name="monotreme-api-v1-Shortcut-OpenGraphMetadata"></a>

### Shortcut.OpenGraphMetadata



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [string](#string) |  |  |
| description | [string](#string) |  |  |
| image | [string](#string) |  |  |






<a name="monotreme-api-v1-UpdateShortcutRequest"></a>

### UpdateShortcutRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| shortcut | [Shortcut](#monotreme-api-v1-Shortcut) |  |  |
| update_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  |  |





 

 

 


<a name="monotreme-api-v1-ShortcutService"></a>

### ShortcutService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListShortcuts | [ListShortcutsRequest](#monotreme-api-v1-ListShortcutsRequest) | [ListShortcutsResponse](#monotreme-api-v1-ListShortcutsResponse) | ListShortcuts returns a list of shortcuts. |
| GetShortcut | [GetShortcutRequest](#monotreme-api-v1-GetShortcutRequest) | [Shortcut](#monotreme-api-v1-Shortcut) | GetShortcut returns a shortcut by id. |
| GetShortcutByName | [GetShortcutByNameRequest](#monotreme-api-v1-GetShortcutByNameRequest) | [Shortcut](#monotreme-api-v1-Shortcut) | GetShortcutByName returns a shortcut by name. |
| CreateShortcut | [CreateShortcutRequest](#monotreme-api-v1-CreateShortcutRequest) | [Shortcut](#monotreme-api-v1-Shortcut) | CreateShortcut creates a shortcut. |
| UpdateShortcut | [UpdateShortcutRequest](#monotreme-api-v1-UpdateShortcutRequest) | [Shortcut](#monotreme-api-v1-Shortcut) | UpdateShortcut updates a shortcut. |
| DeleteShortcut | [DeleteShortcutRequest](#monotreme-api-v1-DeleteShortcutRequest) | [.google.protobuf.Empty](#google-protobuf-Empty) | DeleteShortcut deletes a shortcut by name. |
| GetShortcutAnalytics | [GetShortcutAnalyticsRequest](#monotreme-api-v1-GetShortcutAnalyticsRequest) | [GetShortcutAnalyticsResponse](#monotreme-api-v1-GetShortcutAnalyticsResponse) | GetShortcutAnalytics returns the analytics for a shortcut. |

 



<a name="api_v1_subscription_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v1/subscription_service.proto



<a name="monotreme-api-v1-DeleteSubscriptionRequest"></a>

### DeleteSubscriptionRequest







<a name="monotreme-api-v1-GetSubscriptionRequest"></a>

### GetSubscriptionRequest







<a name="monotreme-api-v1-Subscription"></a>

### Subscription



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| plan | [PlanType](#monotreme-api-v1-PlanType) |  |  |
| started_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| expires_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| features | [string](#string) | repeated |  |
| seats | [int32](#int32) |  |  |
| shortcuts_limit | [int32](#int32) |  |  |
| collections_limit | [int32](#int32) |  |  |






<a name="monotreme-api-v1-UpdateSubscriptionRequest"></a>

### UpdateSubscriptionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| license_key | [string](#string) |  |  |





 


<a name="monotreme-api-v1-PlanType"></a>

### PlanType


| Name | Number | Description |
| ---- | ------ | ----------- |
| PLAN_TYPE_UNSPECIFIED | 0 |  |
| FREE | 1 |  |
| PRO | 2 |  |
| ENTERPRISE | 3 |  |


 

 


<a name="monotreme-api-v1-SubscriptionService"></a>

### SubscriptionService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSubscription | [GetSubscriptionRequest](#monotreme-api-v1-GetSubscriptionRequest) | [Subscription](#monotreme-api-v1-Subscription) | GetSubscription gets the current subscription of Monotreme instance. |
| UpdateSubscription | [UpdateSubscriptionRequest](#monotreme-api-v1-UpdateSubscriptionRequest) | [Subscription](#monotreme-api-v1-Subscription) | UpdateSubscription updates the subscription. |
| DeleteSubscription | [DeleteSubscriptionRequest](#monotreme-api-v1-DeleteSubscriptionRequest) | [Subscription](#monotreme-api-v1-Subscription) | DeleteSubscription deletes the subscription. |

 



<a name="api_v1_user_setting_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v1/user_setting_service.proto



<a name="monotreme-api-v1-GetUserSettingRequest"></a>

### GetUserSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  | id is the user id. |






<a name="monotreme-api-v1-UpdateUserSettingRequest"></a>

### UpdateUserSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  | id is the user id. |
| user_setting | [UserSetting](#monotreme-api-v1-UserSetting) |  | user_setting is the user setting to update. |
| update_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  | update_mask is the field mask to update. |






<a name="monotreme-api-v1-UserSetting"></a>

### UserSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [int32](#int32) |  |  |
| general | [UserSetting.GeneralSetting](#monotreme-api-v1-UserSetting-GeneralSetting) |  |  |
| access_tokens | [UserSetting.AccessTokensSetting](#monotreme-api-v1-UserSetting-AccessTokensSetting) |  |  |






<a name="monotreme-api-v1-UserSetting-AccessTokensSetting"></a>

### UserSetting.AccessTokensSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_tokens | [UserSetting.AccessTokensSetting.AccessToken](#monotreme-api-v1-UserSetting-AccessTokensSetting-AccessToken) | repeated | Nested repeated field |






<a name="monotreme-api-v1-UserSetting-AccessTokensSetting-AccessToken"></a>

### UserSetting.AccessTokensSetting.AccessToken



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  | The access token is a JWT token, including expiration time, issuer, etc. |
| description | [string](#string) |  | A description for the access token. |






<a name="monotreme-api-v1-UserSetting-GeneralSetting"></a>

### UserSetting.GeneralSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| locale | [string](#string) |  |  |
| color_theme | [string](#string) |  |  |





 

 

 


<a name="monotreme-api-v1-UserSettingService"></a>

### UserSettingService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetUserSetting | [GetUserSettingRequest](#monotreme-api-v1-GetUserSettingRequest) | [UserSetting](#monotreme-api-v1-UserSetting) | GetUserSetting returns the user setting. |
| UpdateUserSetting | [UpdateUserSettingRequest](#monotreme-api-v1-UpdateUserSettingRequest) | [UserSetting](#monotreme-api-v1-UserSetting) | UpdateUserSetting updates the user setting. |

 



<a name="api_v1_workspace_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v1/workspace_service.proto



<a name="monotreme-api-v1-GetWorkspaceProfileRequest"></a>

### GetWorkspaceProfileRequest







<a name="monotreme-api-v1-GetWorkspaceSettingRequest"></a>

### GetWorkspaceSettingRequest







<a name="monotreme-api-v1-IdentityProvider"></a>

### IdentityProvider



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The unique identifier of the identity provider. |
| title | [string](#string) |  |  |
| type | [IdentityProvider.Type](#monotreme-api-v1-IdentityProvider-Type) |  |  |
| config | [IdentityProviderConfig](#monotreme-api-v1-IdentityProviderConfig) |  |  |






<a name="monotreme-api-v1-IdentityProviderConfig"></a>

### IdentityProviderConfig



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| oauth2 | [IdentityProviderConfig.OAuth2Config](#monotreme-api-v1-IdentityProviderConfig-OAuth2Config) |  |  |






<a name="monotreme-api-v1-IdentityProviderConfig-FieldMapping"></a>

### IdentityProviderConfig.FieldMapping



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| identifier | [string](#string) |  |  |
| display_name | [string](#string) |  |  |






<a name="monotreme-api-v1-IdentityProviderConfig-OAuth2Config"></a>

### IdentityProviderConfig.OAuth2Config



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| client_id | [string](#string) |  |  |
| client_secret | [string](#string) |  |  |
| auth_url | [string](#string) |  |  |
| token_url | [string](#string) |  |  |
| user_info_url | [string](#string) |  |  |
| scopes | [string](#string) | repeated |  |
| field_mapping | [IdentityProviderConfig.FieldMapping](#monotreme-api-v1-IdentityProviderConfig-FieldMapping) |  |  |






<a name="monotreme-api-v1-UpdateWorkspaceSettingRequest"></a>

### UpdateWorkspaceSettingRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| setting | [WorkspaceSetting](#monotreme-api-v1-WorkspaceSetting) |  | The user setting. |
| update_mask | [google.protobuf.FieldMask](#google-protobuf-FieldMask) |  | The update mask. |






<a name="monotreme-api-v1-WorkspaceProfile"></a>

### WorkspaceProfile



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| mode | [string](#string) |  | Current workspace mode: dev, prod. |
| version | [string](#string) |  | Current workspace version. |
| owner | [string](#string) |  | The owner name. Format: &#34;users/{id}&#34; |
| subscription | [Subscription](#monotreme-api-v1-Subscription) |  | The workspace subscription. |
| custom_style | [string](#string) |  | The custom style. |
| branding | [bytes](#bytes) |  | The workspace branding. |






<a name="monotreme-api-v1-WorkspaceSetting"></a>

### WorkspaceSetting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| instance_url | [string](#string) |  | The url of instance. |
| branding | [bytes](#bytes) |  | The workspace custome branding. |
| custom_style | [string](#string) |  | The custom style. |
| default_visibility | [Visibility](#monotreme-api-v1-Visibility) |  | The default visibility of shortcuts and collections. |
| identity_providers | [IdentityProvider](#monotreme-api-v1-IdentityProvider) | repeated | The identity providers. |
| disallow_user_registration | [bool](#bool) |  | Whether to disallow user registration by email&amp;password. |
| disallow_password_auth | [bool](#bool) |  | Whether to disallow password authentication. |





 


<a name="monotreme-api-v1-IdentityProvider-Type"></a>

### IdentityProvider.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNSPECIFIED | 0 |  |
| OAUTH2 | 1 |  |


 

 


<a name="monotreme-api-v1-WorkspaceService"></a>

### WorkspaceService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetWorkspaceProfile | [GetWorkspaceProfileRequest](#monotreme-api-v1-GetWorkspaceProfileRequest) | [WorkspaceProfile](#monotreme-api-v1-WorkspaceProfile) |  |
| GetWorkspaceSetting | [GetWorkspaceSettingRequest](#monotreme-api-v1-GetWorkspaceSettingRequest) | [WorkspaceSetting](#monotreme-api-v1-WorkspaceSetting) |  |
| UpdateWorkspaceSetting | [UpdateWorkspaceSettingRequest](#monotreme-api-v1-UpdateWorkspaceSettingRequest) | [WorkspaceSetting](#monotreme-api-v1-WorkspaceSetting) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

