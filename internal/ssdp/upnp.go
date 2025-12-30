package ssdp

import (
	"context"
	"fmt"

	"github.com/koron/go-ssdp"
	"github.com/tr1v3r/rcast/internal/upnp"
)

type UpnpAds struct {
	rootAds []*ssdp.Advertiser
}

func NewUpnpAds(ctx context.Context, location, deviceUUID, serverName string, maxAge int) (*UpnpAds, error) {
	var (
		err  error
		usns = []struct{ st, usn string }{
			{upnp.DeviceType, deviceUUID + "::" + upnp.DeviceType},
			{upnp.AVTransportType, deviceUUID + "::" + upnp.AVTransportType},
			{upnp.RenderingType, deviceUUID + "::" + upnp.RenderingType},
			{"upnp:rootdevice", deviceUUID + "::upnp:rootdevice"},
			{st: deviceUUID, usn: deviceUUID},
			{upnp.ContentDirectoryType, deviceUUID},
		}
	)

	advertiserList := make([]*ssdp.Advertiser, len(usns))

	for i, x := range usns {
		advertiserList[i], err = ssdp.Advertise(x.st, x.usn, location, serverName, maxAge)
		if err != nil {
			return nil, fmt.Errorf("advertise error: %w", err)
		}
	}

	return &UpnpAds{
		rootAds: advertiserList,
	}, nil

}

func (u *UpnpAds) NotifyAll() error {
	for _, v := range u.rootAds {
		if err := v.Alive(); err != nil {
			return err
		}
	}
	return nil
}

func (u *UpnpAds) CloseAll() {
	for _, v := range u.rootAds {
		v.Bye()
		v.Close()
	}
}
