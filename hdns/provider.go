package hdns

import (
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"gitlab.com/alxrem/hdns-go/hdns"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HDNS_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"hdns_zone":   resourceZone(),
			"hdns_record": resourceRecord(),
		},
		ConfigureFunc: configureProvider,
	}
}

type providerConfig struct {
	client *hdns.Client
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	opts := []hdns.ClientOption{
		hdns.WithToken(d.Get("token").(string)),
		hdns.WithApplication("hdns-terraform", "0.0.1"),
	}
	if logging.LogLevel() != "" {
		writer, err := logging.LogOutput()
		if err != nil {
			return nil, err
		}
		opts = append(opts, hdns.WithDebugWriter(writer.(*logutils.LevelFilter).Writer))
	}

	client := hdns.NewClient(opts...)

	return &providerConfig{
		client: client,
	}, nil
}
