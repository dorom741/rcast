package ssdp

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/koron/go-ssdp"
	"github.com/tr1v3r/pkg/log"
)

func SetInterface(nic ...net.Interface) { ssdp.Interfaces = nic }

func StartSSDP(ctx context.Context, baseURL, deviceUUID, serverName string) {
	var (
		err      error
		location = fmt.Sprintf("%s/device.xml", baseURL)
		maxAge   = 1800
	)

	ads, err := NewUpnpAds(ctx, location, deviceUUID, serverName, maxAge)
	if err != nil {
		log.CtxError(ctx, "NewUpnpAds error: %s", err)
	}

	// ads, err := ssdp.Advertise(upnp.ContentDirectoryType, deviceUUID, location, serverName, maxAge)
	// if err != nil {
	// 	log.CtxError(ctx, "Advertise error: %s", err)
	// }

	m := &ssdp.Monitor{
		Search: func(m *ssdp.SearchMessage) {
			log.CtxInfo(ctx, "Search: From=%s Type=%s\n", m.From.String(), m.Type)

			isSearch := strings.Contains(m.Type, "ssdp:all") || strings.Contains(m.Type, "service:ContentDirectory") || strings.Contains(m.Type, "service:ConnectionManager") || strings.Contains(m.Type, "device:MediaServer")
			if !isSearch {
				return
			}
			if err := ads.NotifyAll(); err != nil {
				log.CtxError(ctx, "NotifyAll error: %s", err)
			}

			// ad.Alive()
			// for _, advertiser := range advertiserList {
			// 	advertiser.Alive()
			// }
		},
	}
	m.Start()

	repeat := time.Tick(time.Duration(maxAge) * time.Second)
	for {
		select {
		case <-ctx.Done():
			ads.CloseAll()
			return
		case <-repeat:
			if err := ads.NotifyAll(); err != nil {
				log.CtxError(ctx, "NotifyAll error: %s", err)
			}
		}
	}

}
