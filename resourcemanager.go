package mandala

import (
	"unsafe"

	"git.tideland.biz/goas/loop"
)

type LoadResourceResponse struct {
	// The buffer containing the resource.
	Buffer []byte

	// The error eventually generated by the reading operation.
	Error error
}

// The type of request to send to ResourceManager in order to read the
// resource.
type LoadResourceRequest struct {
	// The path of the resource file, for example
	// "res/drawable/gopher.png". ResourcePath will be prefixed to
	// Filename.
	Filename string

	// Response is a channel for receiving the response from the
	// resource manager.
	Response chan LoadResourceResponse
}

func resourceLoopFunc(activity chan unsafe.Pointer, request chan interface{}) loop.LoopFunc {
	var act unsafe.Pointer
	return func(l loop.Loop) error {
		for {
			select {
			case act = <-activity:
			case untypedRequest := <-request:
				switch req := untypedRequest.(type) {
				case LoadResourceRequest:
					buf, err := loadResource(act, req.Filename)
					req.Response <- LoadResourceResponse{buf, err}
				}
			}
		}
	}
}

// ReadResource reads a resource named filename and send the response
// to the given responseCh channel.
func ReadResource(filename string, responseCh chan LoadResourceResponse) {
	request := LoadResourceRequest{
		Filename: filename,
		Response: responseCh,
	}
	ResourceManager() <- request
}
