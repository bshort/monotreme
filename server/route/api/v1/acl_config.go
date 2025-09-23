package v1

import "strings"

var allowedMethodsWhenUnauthorized = map[string]bool{
	"/monotreme.api.v1.WorkspaceService/GetWorkspaceProfile":  true,
	"/monotreme.api.v1.WorkspaceService/GetWorkspaceSetting":  true,
	"/monotreme.api.v1.AuthService/GetAuthStatus":             true,
	"/monotreme.api.v1.AuthService/SignIn":                    true,
	"/monotreme.api.v1.AuthService/SignInWithSSO":             true,
	"/monotreme.api.v1.AuthService/SignUp":                    true,
	"/monotreme.api.v1.AuthService/SignOut":                   true,
	"/monotreme.api.v1.ShortcutService/GetShortcut":           true,
	"/monotreme.api.v1.ShortcutService/GetShortcutByName":     true,
	"/monotreme.api.v1.CollectionService/GetCollectionByName": true,
}

// isUnauthorizeAllowedMethod returns true if the method is allowed to be called when the user is not authorized.
func isUnauthorizeAllowedMethod(methodName string) bool {
	if strings.HasPrefix(methodName, "/grpc.reflection") {
		return true
	}
	return allowedMethodsWhenUnauthorized[methodName]
}

var allowedMethodsOnlyForAdmin = map[string]bool{
	"/monotreme.api.v1.UserService/CreateUser":                  true,
	"/monotreme.api.v1.UserService/DeleteUser":                  true,
	"/monotreme.api.v1.WorkspaceService/UpdateWorkspaceSetting": true,
	"/monotreme.api.v1.SubscriptionService/UpdateSubscription":  true,
}

// isOnlyForAdminAllowedMethod returns true if the method is allowed to be called only by admin.
func isOnlyForAdminAllowedMethod(methodName string) bool {
	return allowedMethodsOnlyForAdmin[methodName]
}
