package client

import (
	"fmt"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"net/http"
	"github.com/layer-x/layerx-commons/lxerrors"
	"encoding/json"
	"github.com/emc-advanced-dev/unik/pkg/types"
)

type volumes struct {
	unikIP string
}

func (v *volumes) All() ([]*types.Volume, error) {
	resp, body, err := lxhttpclient.Get(v.unikIP, "/volumes", nil)
	if err != nil  {
		return nil, lxerrors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var volumes []*types.Volume
	if err := json.Unmarshal(body, &volumes); err != nil {
		return nil, lxerrors.New(fmt.Sprintf("response body %s did not unmarshal to type []*types.Volume", string(body)), err)
	}
	return volumes, nil
}

func (v *volumes) Get(id string) (*types.Volume, error) {
	resp, body, err := lxhttpclient.Get(v.unikIP, "/volumes/"+id, nil)
	if err != nil  {
		return nil, lxerrors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), nil)
	}
	var volume types.Volume
	if err := json.Unmarshal(body, &volume); err != nil {
		return nil, lxerrors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Volume", string(body)), err)
	}
	return &volume, nil
}

func (v *volumes) Delete(id string) error {
	resp, body, err := lxhttpclient.Delete(v.unikIP, "/volumes/"+id, nil)
	if err != nil  {
		return nil, lxerrors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (v *volumes) Create(name, dataTar, provider string, size int) (*types.Volume, error) {
	query := fmt.Sprintf("?size=%v&provider=%v", size, provider)
	//no data provided
	var (
		resp *http.Response
		body []byte
		err error
	)
	if dataTar == "" {
		resp, body, err = lxhttpclient.Post(v.unikIP, "/volumes/" + name + query, nil, nil)
		if err != nil {
			return nil, lxerrors.New("request failed", err)
		}
		if resp.StatusCode != http.StatusCreated {
			return nil, lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
		}
	} else {
		resp, body, err = lxhttpclient.PostFile(v.unikIP, "/volumes/" + name + query, "tarfile", dataTar)
		if err != nil {
			return nil, lxerrors.New("request failed", err)
		}
		if resp.StatusCode != http.StatusCreated {
			return nil, lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
		}
	}
	var volume types.Volume
	if err := json.Unmarshal(body, &volume); err != nil {
		return nil, lxerrors.New(fmt.Sprintf("response body %s did not unmarshal to type *types.Volume", string(body)), err)
	}
	return &volume, nil
}

func (v *volumes) Attach(id, instanceId, mountPoint string) error {
	query := fmt.Sprintf("?mount=%v", mountPoint)
	resp, body, err := lxhttpclient.Post(v.unikIP, "/volumes/"+id+"/attach/"+instanceId+query, nil, nil)
	if err != nil  {
		return nil, lxerrors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}

func (v *volumes) Detach(id string) error {
	resp, body, err := lxhttpclient.Post(v.unikIP, "/volumes/"+id+"/detach", nil, nil)
	if err != nil  {
		return nil, lxerrors.New("request failed", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		return lxerrors.New(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)), err)
	}
	return nil
}