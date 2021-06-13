/*
 * Roster Servise API
 *
 * This is TEST API for my service
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Roster struct {
	// ID roster
	Id string `json:"id"`
	// NAME roster
	Name string `json:"name"`
	// ID user
	IdUser string `json:"idUser"`
	// status roster. 0 - valid, 1 - need update
	Status int32 `json:"status"`

	Units []Unit `json:"units"`
}