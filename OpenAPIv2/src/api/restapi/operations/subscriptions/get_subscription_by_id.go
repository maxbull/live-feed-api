// Code generated by go-swagger; DO NOT EDIT.

package subscriptions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetSubscriptionByIDHandlerFunc turns a function with the right signature into a get subscription by Id handler
type GetSubscriptionByIDHandlerFunc func(GetSubscriptionByIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetSubscriptionByIDHandlerFunc) Handle(params GetSubscriptionByIDParams) middleware.Responder {
	return fn(params)
}

// GetSubscriptionByIDHandler interface for that can handle valid get subscription by Id params
type GetSubscriptionByIDHandler interface {
	Handle(GetSubscriptionByIDParams) middleware.Responder
}

// NewGetSubscriptionByID creates a new http.Handler for the get subscription by Id operation
func NewGetSubscriptionByID(ctx *middleware.Context, handler GetSubscriptionByIDHandler) *GetSubscriptionByID {
	return &GetSubscriptionByID{Context: ctx, Handler: handler}
}

/*GetSubscriptionByID swagger:route GET /subscriptions/{subscriptionId} subscriptions getSubscriptionById

Info for a specific subscription

*/
type GetSubscriptionByID struct {
	Context *middleware.Context
	Handler GetSubscriptionByIDHandler
}

func (o *GetSubscriptionByID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetSubscriptionByIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
