package api

import (
	"encoding/json"
	"jwt_auth.com/internal/auth"
	"net/http"
)

type Handle struct {
	isAuth bool
	path   string
	handle func(r *http.Request) ([]byte, int)
}

func newHandle(auth bool, path string, handle func(r *http.Request) ([]byte, int)) *Handle {
	return &Handle{
		isAuth: auth,
		path:   path,
		handle: handle,
	}
}

func (h *Handle) Handle(w http.ResponseWriter, r *http.Request) {
	if h.isAuth {
		token := r.Header.Get("Authorization")

		res, err := auth.AuthCheck(token)
		if err != nil {
			json, _ := json.Marshal("err: " + err.Error())
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			return
		}

		if !res {
			json, _ := json.Marshal("err: ваш токен авторизации неверен или просрочен")
			w.WriteHeader(http.StatusNotAcceptable)
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			return
		}
	}

	json, status := h.handle(r)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func RegisterMux() *http.ServeMux {
	mux := http.ServeMux{}
	for path, handle := range urls {
		mux.HandleFunc(path, handle.Handle)
	}
	return &mux
}
