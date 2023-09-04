package utils

import (
	"crypto/tls"
	"crypto/x509"
	"time"

	"main/packages/apiGO"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/bot/basic"
	"github.com/Tnze/go-mc/bot/msg"
	"github.com/Tnze/go-mc/bot/screen"
	"github.com/Tnze/go-mc/bot/world"
)

var (
	Roots                   *x509.CertPool = x509.NewCertPool()
	Con                     Config
	Proxy                   apiGO.Proxys
	Bearer                  apiGO.MCbearers
	RGB                     []string
	First_mfa               bool = true
	First_gc                bool = true
	Use_gc, Use_mfa, Accamt int
	Accs                    map[string][]Proxys_Accs = make(map[string][]Proxys_Accs)
	client                  *bot.Client
	player                  *basic.Player
	chatHandler             *msg.Manager
	worldManager            *world.World
	screenManager           *screen.Manager
	letterRunes             = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	ProxyByte               = []byte(`
	-- GlobalSign Root R2, valid until Dec 15, 2021
	-----BEGIN CERTIFICATE-----
	MIIDujCCAqKgAwIBAgILBAAAAAABD4Ym5g0wDQYJKoZIhvcNAQEFBQAwTDEgMB4G
	A1UECxMXR2xvYmFsU2lnbiBSb290IENBIC0gUjIxEzARBgNVBAoTCkdsb2JhbFNp
	Z24xEzARBgNVBAMTCkdsb2JhbFNpZ24wHhcNMDYxMjE1MDgwMDAwWhcNMjExMjE1
	MDgwMDAwWjBMMSAwHgYDVQQLExdHbG9iYWxTaWduIFJvb3QgQ0EgLSBSMjETMBEG
	A1UEChMKR2xvYmFsU2lnbjETMBEGA1UEAxMKR2xvYmFsU2lnbjCCASIwDQYJKoZI
	hvcNAQEBBQADggEPADCCAQoCggEBAKbPJA6+Lm8omUVCxKs+IVSbC9N/hHD6ErPL
	v4dfxn+G07IwXNb9rfF73OX4YJYJkhD10FPe+3t+c4isUoh7SqbKSaZeqKeMWhG8
	eoLrvozps6yWJQeXSpkqBy+0Hne/ig+1AnwblrjFuTosvNYSuetZfeLQBoZfXklq
	tTleiDTsvHgMCJiEbKjNS7SgfQx5TfC4LcshytVsW33hoCmEofnTlEnLJGKRILzd
	C9XZzPnqJworc5HGnRusyMvo4KD0L5CLTfuwNhv2GXqF4G3yYROIXJ/gkwpRl4pa
	zq+r1feqCapgvdzZX99yqWATXgAByUr6P6TqBwMhAo6CygPCm48CAwEAAaOBnDCB
	mTAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUm+IH
	V2ccHsBqBt5ZtJot39wZhi4wNgYDVR0fBC8wLTAroCmgJ4YlaHR0cDovL2NybC5n
	bG9iYWxzaWduLm5ldC9yb290LXIyLmNybDAfBgNVHSMEGDAWgBSb4gdXZxwewGoG
	3lm0mi3f3BmGLjANBgkqhkiG9w0BAQUFAAOCAQEAmYFThxxol4aR7OBKuEQLq4Gs
	J0/WwbgcQ3izDJr86iw8bmEbTUsp9Z8FHSbBuOmDAGJFtqkIk7mpM0sYmsL4h4hO
	291xNBrBVNpGP+DTKqttVCL1OmLNIG+6KYnX3ZHu01yiPqFbQfXf5WRDLenVOavS
	ot+3i9DAgBkcRcAtjOj4LaR0VknFBbVPFd5uRHg5h6h+u/N5GJG79G+dwfCMNYxd
	AfvDbbnvRG15RjF+Cv6pgsH/76tuIMRQyV+dTZsXjAzlAcmgQWpzU/qlULRuJQ/7
	TBj0/VLZjmmx6BEP3ojY+x1J96relc8geMJgEtslQIxq/H5COEBkEveegeGTLg==
	-----END CERTIFICATE-----`)
)

type SniperProxy struct {
	Proxy        *tls.Conn
	UsedAt       time.Time
	Alive        bool
	ProxyDetails Proxies
}

type Proxys_Accs struct {
	Proxy string
	Accs  []apiGO.Info
}

type NameMCInfo struct {
	Action string     `json:"action"`
	Desc   string     `json:"desc"`
	Code   string     `json:"code"`
	Data   NameMCData `json:"data"`
}
type NameMCData struct {
	Status    string    `json:"status"`
	Searches  string    `json:"searches"`
	Begin     int       `json:"begin"`
	End       int       `json:"end"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Headurl   string    `json:"headurl"`
}

type NameMCHead struct {
	Bodyurl string `json:"bodyurl"`
	Headurl string `json:"headurl"`
	ID      string `json:"id"`
}

type Names struct {
	Name  string
	Taken bool
}

type Proxies struct {
	IP, Port, User, Password string
}

type Status struct {
	Data struct {
		Status string `json:"status"`
	} `json:"details"`
}

type CF struct {
	Tokens   string `json:"tokens"`
	GennedAT int64  `json:"unix_of_creation"`
}

type Config struct {
	Gradient   []Values        `json:"gradient"`
	NMC        Namemc_Data     `json:"namemc_settings"`
	Settings   AccountSettings `json:"settings"`
	Bools      Bools           `json:"sniper_config"`
	SkinChange Skin            `json:"skin_config"`
	CF         CF              `json:"cf_tokens"`
	Bearers    []Bearers       `json:"Bearers"`
	Recovery   []Succesful     `json:"recovery"`
}

type Namemc_Data struct {
	UseNMC          bool       `json:"usenamemc_fordroptime_andautofollow"`
	Display         string     `json:"name_to_use_for_follows"`
	Key             string     `json:"namemc_email:pass"`
	NamemcLoginData NMC        `json:"namemc_login_data"`
	P               []Profiles `json:"genned_profiles"`
}

type Profiles struct {
	Session_ID string `json:"session_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type NMC struct {
	Token      string `json:"token"`
	LastAuthed int64  `json:"last_unix_auth_timestamp"`
}

type Bearers struct {
	Bearer               string   `json:"Bearer"`
	Email                string   `json:"Email"`
	Password             string   `json:"Password"`
	AuthInterval         int64    `json:"AuthInterval"`
	AuthedAt             int64    `json:"AuthedAt"`
	Type                 string   `json:"Type"`
	NameChange           bool     `json:"NameChange"`
	Info                 UserINFO `json:"Info"`
	NOT_ENTITLED_CHECKED bool     `json:"checked_entitled"`
}

type Succesful struct {
	Email     string
	Recovery  string
	Code_Used string
}
type Data struct {
	Info []Succesful
}

type Refresh struct {
	Time_since_last_gen int64 `json:"last_entitled_prevention"`
}

type AccountSettings struct {
	Youtube          string `json:"youtube_link"`
	AskForUnixPrompt bool   `json:"ask_for_unix_prompt"`
	AccountsPerGc    int    `json:"accounts_per_gc_proxy"`
	AccountsPerMfa   int    `json:"accounts_per_mfa_proxy"`
	GC_ReqAmt        int    `json:"amt_reqs_per_gc_acc"`
	MFA_ReqAmt       int    `json:"amt_reqs_per_mfa_acc"`
	SleepAmtPerGc    int    `json:"sleep_for_gc"`
	SleepAmtPerMfa   int    `json:"sleep_for_mfa"`
	UseCustomSpread  bool   `json:"use_own_spread_value"`
	Spread           int64  `json:"spread_ms"`
}

type Bools struct {
	UseCF                            bool `json:"use_cf_token"`
	UseProxyDuringAuth               bool `json:"useproxysduringauth"`
	UseWebhook                       bool `json:"sendpersonalwhonsnipe"`
	FirstUse                         bool `json:"firstuse_IGNORETHIS"`
	DownloadedPW                     bool `json:"pwinstalled_IGNORETHIS"`
	ApplyNewRecoveryToExistingEmails bool `json:"applynewemailstoexistingrecoveryemails"`
}

type Values struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

type Info struct {
	Bearer       string
	RefreshToken string
	AccessToken  string
	Expires      int
	AccountType  string
	Email        string
	Password     string
	Requests     int
	Info         UserINFO `json:"Info"`
	Error        string
}

type UserINFO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Skin struct {
	Link    string `json:"url"`
	Variant string `json:"variant"`
}

type Payload_auth struct {
	Proxy    string
	Accounts []string
}

type YtPL struct {
	ResponseContext struct {
		ServiceTrackingParams []struct {
			Service string `json:"service"`
			Params  []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"params"`
		} `json:"serviceTrackingParams"`
		MainAppWebResponseContext struct {
			DatasyncID    string `json:"datasyncId"`
			LoggedOut     bool   `json:"loggedOut"`
			TrackingParam string `json:"trackingParam"`
		} `json:"mainAppWebResponseContext"`
		WebResponseContextExtensionData struct {
			HasDecorated bool `json:"hasDecorated"`
		} `json:"webResponseContextExtensionData"`
	} `json:"responseContext"`
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Selected       bool   `json:"selected"`
					TrackingParams string `json:"trackingParams"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
	Alerts []struct {
		AlertWithButtonRenderer struct {
			Type string `json:"type"`
			Text struct {
				SimpleText string `json:"simpleText"`
			} `json:"text"`
			DismissButton struct {
				ButtonRenderer struct {
					Style      string `json:"style"`
					Size       string `json:"size"`
					IsDisabled bool   `json:"isDisabled"`
					Icon       struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					TrackingParams    string `json:"trackingParams"`
					AccessibilityData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibilityData"`
				} `json:"buttonRenderer"`
			} `json:"dismissButton"`
		} `json:"alertWithButtonRenderer"`
	} `json:"alerts"`
	Metadata struct {
		PlaylistMetadataRenderer struct {
			Title                  string `json:"title"`
			AndroidAppindexingLink string `json:"androidAppindexingLink"`
			IosAppindexingLink     string `json:"iosAppindexingLink"`
		} `json:"playlistMetadataRenderer"`
	} `json:"metadata"`
	TrackingParams string `json:"trackingParams"`
	Microformat    struct {
		MicroformatDataRenderer struct {
			URLCanonical string `json:"urlCanonical"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			Thumbnail    struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			SiteName           string `json:"siteName"`
			AppName            string `json:"appName"`
			AndroidPackage     string `json:"androidPackage"`
			IosAppStoreID      string `json:"iosAppStoreId"`
			IosAppArguments    string `json:"iosAppArguments"`
			OgType             string `json:"ogType"`
			URLApplinksWeb     string `json:"urlApplinksWeb"`
			URLApplinksIos     string `json:"urlApplinksIos"`
			URLApplinksAndroid string `json:"urlApplinksAndroid"`
			URLTwitterIos      string `json:"urlTwitterIos"`
			URLTwitterAndroid  string `json:"urlTwitterAndroid"`
			TwitterCardType    string `json:"twitterCardType"`
			TwitterSiteHandle  string `json:"twitterSiteHandle"`
			SchemaDotOrgType   string `json:"schemaDotOrgType"`
			Noindex            bool   `json:"noindex"`
			Unlisted           bool   `json:"unlisted"`
			LinkAlternates     []struct {
				HrefURL string `json:"hrefUrl"`
			} `json:"linkAlternates"`
		} `json:"microformatDataRenderer"`
	} `json:"microformat"`
	OnResponseReceivedActions []struct {
		ClickTrackingParams           string `json:"clickTrackingParams"`
		AppendContinuationItemsAction struct {
			ContinuationItems []struct {
				PlaylistVideoRenderer struct {
					VideoID   string `json:"videoId"`
					Thumbnail struct {
						Thumbnails []struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"thumbnails"`
					} `json:"thumbnail"`
					Title struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
						Accessibility struct {
							AccessibilityData struct {
								Label string `json:"label"`
							} `json:"accessibilityData"`
						} `json:"accessibility"`
					} `json:"title"`
					Index struct {
						SimpleText string `json:"simpleText"`
					} `json:"index"`
					ShortBylineText struct {
						Runs []struct {
							Text               string `json:"text"`
							NavigationEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										URL         string `json:"url"`
										WebPageType string `json:"webPageType"`
										RootVe      int    `json:"rootVe"`
										APIURL      string `json:"apiUrl"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								BrowseEndpoint struct {
									BrowseID         string `json:"browseId"`
									CanonicalBaseURL string `json:"canonicalBaseUrl"`
								} `json:"browseEndpoint"`
							} `json:"navigationEndpoint"`
						} `json:"runs"`
					} `json:"shortBylineText"`
					LengthText struct {
						Accessibility struct {
							AccessibilityData struct {
								Label string `json:"label"`
							} `json:"accessibilityData"`
						} `json:"accessibility"`
						SimpleText string `json:"simpleText"`
					} `json:"lengthText"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						WatchEndpoint struct {
							VideoID        string `json:"videoId"`
							PlaylistID     string `json:"playlistId"`
							Index          int    `json:"index"`
							Params         string `json:"params"`
							PlayerParams   string `json:"playerParams"`
							LoggingContext struct {
								VssLoggingContext struct {
									SerializedContextData string `json:"serializedContextData"`
								} `json:"vssLoggingContext"`
							} `json:"loggingContext"`
							WatchEndpointSupportedOnesieConfig struct {
								HTML5PlaybackOnesieConfig struct {
									CommonConfig struct {
										URL string `json:"url"`
									} `json:"commonConfig"`
								} `json:"html5PlaybackOnesieConfig"`
							} `json:"watchEndpointSupportedOnesieConfig"`
						} `json:"watchEndpoint"`
					} `json:"navigationEndpoint"`
					SetVideoID     string `json:"setVideoId"`
					LengthSeconds  string `json:"lengthSeconds"`
					TrackingParams string `json:"trackingParams"`
					IsPlayable     bool   `json:"isPlayable"`
					Menu           struct {
						MenuRenderer struct {
							Items []struct {
								MenuServiceItemRenderer struct {
									Text struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"text"`
									Icon struct {
										IconType string `json:"iconType"`
									} `json:"icon"`
									ServiceEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												SendPost bool `json:"sendPost"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										SignalServiceEndpoint struct {
											Signal  string `json:"signal"`
											Actions []struct {
												ClickTrackingParams  string `json:"clickTrackingParams"`
												AddToPlaylistCommand struct {
													OpenMiniplayer      bool   `json:"openMiniplayer"`
													VideoID             string `json:"videoId"`
													ListType            string `json:"listType"`
													OnCreateListCommand struct {
														ClickTrackingParams string `json:"clickTrackingParams"`
														CommandMetadata     struct {
															WebCommandMetadata struct {
																SendPost bool   `json:"sendPost"`
																APIURL   string `json:"apiUrl"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
														CreatePlaylistServiceEndpoint struct {
															VideoIds []string `json:"videoIds"`
															Params   string   `json:"params"`
														} `json:"createPlaylistServiceEndpoint"`
													} `json:"onCreateListCommand"`
													VideoIds []string `json:"videoIds"`
												} `json:"addToPlaylistCommand"`
											} `json:"actions"`
										} `json:"signalServiceEndpoint"`
									} `json:"serviceEndpoint"`
									TrackingParams string `json:"trackingParams"`
								} `json:"menuServiceItemRenderer,omitempty"`
								MenuServiceItemDownloadRenderer struct {
									ServiceEndpoint struct {
										ClickTrackingParams  string `json:"clickTrackingParams"`
										OfflineVideoEndpoint struct {
											VideoID      string `json:"videoId"`
											OnAddCommand struct {
												ClickTrackingParams      string `json:"clickTrackingParams"`
												GetDownloadActionCommand struct {
													VideoID string `json:"videoId"`
													Params  string `json:"params"`
												} `json:"getDownloadActionCommand"`
											} `json:"onAddCommand"`
										} `json:"offlineVideoEndpoint"`
									} `json:"serviceEndpoint"`
									TrackingParams string `json:"trackingParams"`
								} `json:"menuServiceItemDownloadRenderer,omitempty"`
							} `json:"items"`
							TrackingParams string `json:"trackingParams"`
							Accessibility  struct {
								AccessibilityData struct {
									Label string `json:"label"`
								} `json:"accessibilityData"`
							} `json:"accessibility"`
						} `json:"menuRenderer"`
					} `json:"menu"`
					ThumbnailOverlays []struct {
						ThumbnailOverlayPlaybackStatusRenderer struct {
							Texts []struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"texts"`
						} `json:"thumbnailOverlayPlaybackStatusRenderer,omitempty"`
						ThumbnailOverlayResumePlaybackRenderer struct {
							PercentDurationWatched int `json:"percentDurationWatched"`
						} `json:"thumbnailOverlayResumePlaybackRenderer,omitempty"`
						ThumbnailOverlayTimeStatusRenderer struct {
							Text struct {
								Accessibility struct {
									AccessibilityData struct {
										Label string `json:"label"`
									} `json:"accessibilityData"`
								} `json:"accessibility"`
								SimpleText string `json:"simpleText"`
							} `json:"text"`
							Style string `json:"style"`
						} `json:"thumbnailOverlayTimeStatusRenderer,omitempty"`
						ThumbnailOverlayNowPlayingRenderer struct {
							Text struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"text"`
						} `json:"thumbnailOverlayNowPlayingRenderer,omitempty"`
					} `json:"thumbnailOverlays"`
					VideoInfo struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"videoInfo"`
				} `json:"playlistVideoRenderer"`
			} `json:"continuationItems"`
			TargetID string `json:"targetId"`
		} `json:"appendContinuationItemsAction"`
	} `json:"onResponseReceivedActions"`
	Sidebar struct {
		PlaylistSidebarRenderer struct {
			Items []struct {
				PlaylistSidebarPrimaryInfoRenderer struct {
					ThumbnailRenderer struct {
						PlaylistVideoThumbnailRenderer struct {
							Thumbnail struct {
								Thumbnails []struct {
									URL    string `json:"url"`
									Width  int    `json:"width"`
									Height int    `json:"height"`
								} `json:"thumbnails"`
							} `json:"thumbnail"`
							TrackingParams string `json:"trackingParams"`
						} `json:"playlistVideoThumbnailRenderer"`
					} `json:"thumbnailRenderer"`
					Stats []struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs,omitempty"`
						SimpleText string `json:"simpleText,omitempty"`
					} `json:"stats"`
					Menu struct {
						MenuRenderer struct {
							Items []struct {
								MenuServiceItemRenderer struct {
									Text struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"text"`
									Icon struct {
										IconType string `json:"iconType"`
									} `json:"icon"`
									ServiceEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												SendPost bool `json:"sendPost"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										SignalServiceEndpoint struct {
											Signal  string `json:"signal"`
											Actions []struct {
												ClickTrackingParams             string `json:"clickTrackingParams"`
												OpenOnePickAddVideoModalCommand struct {
													ListID            string `json:"listId"`
													ModalTitle        string `json:"modalTitle"`
													SelectButtonLabel string `json:"selectButtonLabel"`
												} `json:"openOnePickAddVideoModalCommand"`
											} `json:"actions"`
										} `json:"signalServiceEndpoint"`
									} `json:"serviceEndpoint"`
									TrackingParams string `json:"trackingParams"`
								} `json:"menuServiceItemRenderer,omitempty"`
								MenuNavigationItemRenderer struct {
									Text struct {
										SimpleText string `json:"simpleText"`
									} `json:"text"`
									Icon struct {
										IconType string `json:"iconType"`
									} `json:"icon"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												URL         string `json:"url"`
												WebPageType string `json:"webPageType"`
												RootVe      int    `json:"rootVe"`
												APIURL      string `json:"apiUrl"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										BrowseEndpoint struct {
											BrowseID       string `json:"browseId"`
											Params         string `json:"params"`
											Nofollow       bool   `json:"nofollow"`
											NavigationType string `json:"navigationType"`
										} `json:"browseEndpoint"`
									} `json:"navigationEndpoint"`
									TrackingParams string `json:"trackingParams"`
								} `json:"menuNavigationItemRenderer,omitempty"`
							} `json:"items"`
							TrackingParams  string `json:"trackingParams"`
							TopLevelButtons []struct {
								ButtonRenderer struct {
									Style      string `json:"style"`
									Size       string `json:"size"`
									IsDisabled bool   `json:"isDisabled"`
									Icon       struct {
										IconType string `json:"iconType"`
									} `json:"icon"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												URL         string `json:"url"`
												WebPageType string `json:"webPageType"`
												RootVe      int    `json:"rootVe"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										WatchEndpoint struct {
											VideoID        string `json:"videoId"`
											PlaylistID     string `json:"playlistId"`
											Params         string `json:"params"`
											PlayerParams   string `json:"playerParams"`
											LoggingContext struct {
												VssLoggingContext struct {
													SerializedContextData string `json:"serializedContextData"`
												} `json:"vssLoggingContext"`
											} `json:"loggingContext"`
											WatchEndpointSupportedOnesieConfig struct {
												HTML5PlaybackOnesieConfig struct {
													CommonConfig struct {
														URL string `json:"url"`
													} `json:"commonConfig"`
												} `json:"html5PlaybackOnesieConfig"`
											} `json:"watchEndpointSupportedOnesieConfig"`
										} `json:"watchEndpoint"`
									} `json:"navigationEndpoint"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									Tooltip        string `json:"tooltip"`
									TrackingParams string `json:"trackingParams"`
								} `json:"buttonRenderer,omitempty"`
							} `json:"topLevelButtons"`
							Accessibility struct {
								AccessibilityData struct {
									Label string `json:"label"`
								} `json:"accessibilityData"`
							} `json:"accessibility"`
							TargetID string `json:"targetId"`
						} `json:"menuRenderer"`
					} `json:"menu"`
					ThumbnailOverlays []struct {
						ThumbnailOverlaySidePanelRenderer struct {
							Text struct {
								SimpleText string `json:"simpleText"`
							} `json:"text"`
							Icon struct {
								IconType string `json:"iconType"`
							} `json:"icon"`
						} `json:"thumbnailOverlaySidePanelRenderer"`
					} `json:"thumbnailOverlays"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						WatchEndpoint struct {
							VideoID        string `json:"videoId"`
							PlaylistID     string `json:"playlistId"`
							PlayerParams   string `json:"playerParams"`
							LoggingContext struct {
								VssLoggingContext struct {
									SerializedContextData string `json:"serializedContextData"`
								} `json:"vssLoggingContext"`
							} `json:"loggingContext"`
							WatchEndpointSupportedOnesieConfig struct {
								HTML5PlaybackOnesieConfig struct {
									CommonConfig struct {
										URL string `json:"url"`
									} `json:"commonConfig"`
								} `json:"html5PlaybackOnesieConfig"`
							} `json:"watchEndpointSupportedOnesieConfig"`
						} `json:"watchEndpoint"`
					} `json:"navigationEndpoint"`
					ShowMoreText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"showMoreText"`
					TitleForm struct {
						InlineFormRenderer struct {
							FormField struct {
								TextInputFormFieldRenderer struct {
									Label struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"label"`
									Value             string `json:"value"`
									MaxCharacterLimit int    `json:"maxCharacterLimit"`
									Key               string `json:"key"`
									OnChange          struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												SendPost bool   `json:"sendPost"`
												APIURL   string `json:"apiUrl"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										PlaylistEditEndpoint struct {
											PlaylistID string `json:"playlistId"`
											Actions    []struct {
												Action       string `json:"action"`
												PlaylistName string `json:"playlistName"`
											} `json:"actions"`
										} `json:"playlistEditEndpoint"`
									} `json:"onChange"`
									ValidValueRegexp         string `json:"validValueRegexp"`
									InvalidValueErrorMessage struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"invalidValueErrorMessage"`
									Required bool `json:"required"`
								} `json:"textInputFormFieldRenderer"`
							} `json:"formField"`
							EditButton struct {
								ButtonRenderer struct {
									Style      string `json:"style"`
									Size       string `json:"size"`
									IsDisabled bool   `json:"isDisabled"`
									Icon       struct {
										IconType string `json:"iconType"`
									} `json:"icon"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									Tooltip        string `json:"tooltip"`
									TrackingParams string `json:"trackingParams"`
								} `json:"buttonRenderer"`
							} `json:"editButton"`
							SaveButton struct {
								ButtonRenderer struct {
									Style      string `json:"style"`
									Size       string `json:"size"`
									IsDisabled bool   `json:"isDisabled"`
									Text       struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"text"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									TrackingParams string `json:"trackingParams"`
								} `json:"buttonRenderer"`
							} `json:"saveButton"`
							CancelButton struct {
								ButtonRenderer struct {
									Style      string `json:"style"`
									Size       string `json:"size"`
									IsDisabled bool   `json:"isDisabled"`
									Text       struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"text"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									TrackingParams string `json:"trackingParams"`
								} `json:"buttonRenderer"`
							} `json:"cancelButton"`
							TextDisplayed struct {
								SimpleText string `json:"simpleText"`
							} `json:"textDisplayed"`
							Style string `json:"style"`
						} `json:"inlineFormRenderer"`
					} `json:"titleForm"`
					DescriptionForm struct {
						InlineFormRenderer struct {
							FormField struct {
								TextInputFormFieldRenderer struct {
									Label struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"label"`
									Value             string `json:"value"`
									MaxCharacterLimit int    `json:"maxCharacterLimit"`
									Key               string `json:"key"`
									OnChange          struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												SendPost bool   `json:"sendPost"`
												APIURL   string `json:"apiUrl"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										PlaylistEditEndpoint struct {
											PlaylistID string `json:"playlistId"`
											Actions    []struct {
												Action              string `json:"action"`
												PlaylistDescription string `json:"playlistDescription"`
											} `json:"actions"`
											Params string `json:"params"`
										} `json:"playlistEditEndpoint"`
									} `json:"onChange"`
									ValidValueRegexp         string `json:"validValueRegexp"`
									InvalidValueErrorMessage struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"invalidValueErrorMessage"`
									IsMultiline bool `json:"isMultiline"`
								} `json:"textInputFormFieldRenderer"`
							} `json:"formField"`
							EditButton struct {
								ButtonRenderer struct {
									Style      string `json:"style"`
									Size       string `json:"size"`
									IsDisabled bool   `json:"isDisabled"`
									Icon       struct {
										IconType string `json:"iconType"`
									} `json:"icon"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									Tooltip        string `json:"tooltip"`
									TrackingParams string `json:"trackingParams"`
								} `json:"buttonRenderer"`
							} `json:"editButton"`
							SaveButton struct {
								ButtonRenderer struct {
									Style      string `json:"style"`
									Size       string `json:"size"`
									IsDisabled bool   `json:"isDisabled"`
									Text       struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"text"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									TrackingParams string `json:"trackingParams"`
								} `json:"buttonRenderer"`
							} `json:"saveButton"`
							CancelButton struct {
								ButtonRenderer struct {
									Style      string `json:"style"`
									Size       string `json:"size"`
									IsDisabled bool   `json:"isDisabled"`
									Text       struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"text"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									TrackingParams string `json:"trackingParams"`
								} `json:"buttonRenderer"`
							} `json:"cancelButton"`
							Style       string `json:"style"`
							Placeholder struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"placeholder"`
						} `json:"inlineFormRenderer"`
					} `json:"descriptionForm"`
					PrivacyForm struct {
						DropdownFormFieldRenderer struct {
							Dropdown struct {
								DropdownRenderer struct {
									Entries []struct {
										PrivacyDropdownItemRenderer struct {
											Label struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"label"`
											Icon struct {
												IconType string `json:"iconType"`
											} `json:"icon"`
											Description struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"description"`
											Int32Value    int  `json:"int32Value"`
											IsSelected    bool `json:"isSelected"`
											Accessibility struct {
												Label string `json:"label"`
											} `json:"accessibility"`
										} `json:"privacyDropdownItemRenderer"`
									} `json:"entries"`
								} `json:"dropdownRenderer"`
							} `json:"dropdown"`
							Key      string `json:"key"`
							OnChange struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										SendPost bool   `json:"sendPost"`
										APIURL   string `json:"apiUrl"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								PlaylistEditEndpoint struct {
									PlaylistID string `json:"playlistId"`
									Actions    []struct {
										Action          string `json:"action"`
										PlaylistPrivacy string `json:"playlistPrivacy"`
									} `json:"actions"`
									Params string `json:"params"`
								} `json:"playlistEditEndpoint"`
							} `json:"onChange"`
						} `json:"dropdownFormFieldRenderer"`
					} `json:"privacyForm"`
				} `json:"playlistSidebarPrimaryInfoRenderer,omitempty"`
				PlaylistSidebarSecondaryInfoRenderer struct {
					VideoOwner struct {
						VideoOwnerRenderer struct {
							Thumbnail struct {
								Thumbnails []struct {
									URL    string `json:"url"`
									Width  int    `json:"width"`
									Height int    `json:"height"`
								} `json:"thumbnails"`
							} `json:"thumbnail"`
							Title struct {
								Runs []struct {
									Text               string `json:"text"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												URL         string `json:"url"`
												WebPageType string `json:"webPageType"`
												RootVe      int    `json:"rootVe"`
												APIURL      string `json:"apiUrl"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										BrowseEndpoint struct {
											BrowseID         string `json:"browseId"`
											CanonicalBaseURL string `json:"canonicalBaseUrl"`
										} `json:"browseEndpoint"`
									} `json:"navigationEndpoint"`
								} `json:"runs"`
							} `json:"title"`
							NavigationEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										URL         string `json:"url"`
										WebPageType string `json:"webPageType"`
										RootVe      int    `json:"rootVe"`
										APIURL      string `json:"apiUrl"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								BrowseEndpoint struct {
									BrowseID         string `json:"browseId"`
									CanonicalBaseURL string `json:"canonicalBaseUrl"`
								} `json:"browseEndpoint"`
							} `json:"navigationEndpoint"`
							TrackingParams string `json:"trackingParams"`
						} `json:"videoOwnerRenderer"`
					} `json:"videoOwner"`
				} `json:"playlistSidebarSecondaryInfoRenderer,omitempty"`
			} `json:"items"`
			TrackingParams string `json:"trackingParams"`
		} `json:"playlistSidebarRenderer"`
	} `json:"sidebar"`
}

type YTPageConfig struct {
	ResponseContext struct {
		ServiceTrackingParams []struct {
			Service string `json:"service"`
			Params  []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"params"`
		} `json:"serviceTrackingParams"`
		MainAppWebResponseContext struct {
			LoggedOut     bool   `json:"loggedOut"`
			TrackingParam string `json:"trackingParam"`
		} `json:"mainAppWebResponseContext"`
		WebResponseContextExtensionData struct {
			YtConfigData struct {
				VisitorData           string `json:"visitorData"`
				RootVisualElementType int    `json:"rootVisualElementType"`
			} `json:"ytConfigData"`
			HasDecorated bool `json:"hasDecorated"`
		} `json:"webResponseContextExtensionData"`
	} `json:"responseContext"`
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Selected bool `json:"selected"`
					Content  struct {
						SectionListRenderer struct {
							Contents []struct {
								ItemSectionRenderer struct {
									Contents []struct {
										PlaylistVideoListRenderer struct {
											Contents []struct {
												PlaylistVideoRenderer struct {
													VideoID   string `json:"videoId"`
													Thumbnail struct {
														Thumbnails []struct {
															URL    string `json:"url"`
															Width  int    `json:"width"`
															Height int    `json:"height"`
														} `json:"thumbnails"`
													} `json:"thumbnail"`
													Title struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
														Accessibility struct {
															AccessibilityData struct {
																Label string `json:"label"`
															} `json:"accessibilityData"`
														} `json:"accessibility"`
													} `json:"title"`
													Index struct {
														SimpleText string `json:"simpleText"`
													} `json:"index"`
													ShortBylineText struct {
														Runs []struct {
															Text               string `json:"text"`
															NavigationEndpoint struct {
																ClickTrackingParams string `json:"clickTrackingParams"`
																CommandMetadata     struct {
																	WebCommandMetadata struct {
																		URL         string `json:"url"`
																		WebPageType string `json:"webPageType"`
																		RootVe      int    `json:"rootVe"`
																		APIURL      string `json:"apiUrl"`
																	} `json:"webCommandMetadata"`
																} `json:"commandMetadata"`
																BrowseEndpoint struct {
																	BrowseID         string `json:"browseId"`
																	CanonicalBaseURL string `json:"canonicalBaseUrl"`
																} `json:"browseEndpoint"`
															} `json:"navigationEndpoint"`
														} `json:"runs"`
													} `json:"shortBylineText"`
													LengthText struct {
														Accessibility struct {
															AccessibilityData struct {
																Label string `json:"label"`
															} `json:"accessibilityData"`
														} `json:"accessibility"`
														SimpleText string `json:"simpleText"`
													} `json:"lengthText"`
													NavigationEndpoint struct {
														ClickTrackingParams string `json:"clickTrackingParams"`
														CommandMetadata     struct {
															WebCommandMetadata struct {
																URL         string `json:"url"`
																WebPageType string `json:"webPageType"`
																RootVe      int    `json:"rootVe"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
														WatchEndpoint struct {
															VideoID        string `json:"videoId"`
															PlaylistID     string `json:"playlistId"`
															Index          int    `json:"index"`
															Params         string `json:"params"`
															PlayerParams   string `json:"playerParams"`
															LoggingContext struct {
																VssLoggingContext struct {
																	SerializedContextData string `json:"serializedContextData"`
																} `json:"vssLoggingContext"`
															} `json:"loggingContext"`
															WatchEndpointSupportedOnesieConfig struct {
																HTML5PlaybackOnesieConfig struct {
																	CommonConfig struct {
																		URL string `json:"url"`
																	} `json:"commonConfig"`
																} `json:"html5PlaybackOnesieConfig"`
															} `json:"watchEndpointSupportedOnesieConfig"`
														} `json:"watchEndpoint"`
													} `json:"navigationEndpoint"`
													LengthSeconds  string `json:"lengthSeconds"`
													TrackingParams string `json:"trackingParams"`
													IsPlayable     bool   `json:"isPlayable"`
													Menu           struct {
														MenuRenderer struct {
															Items []struct {
																MenuServiceItemRenderer struct {
																	Text struct {
																		Runs []struct {
																			Text string `json:"text"`
																		} `json:"runs"`
																	} `json:"text"`
																	Icon struct {
																		IconType string `json:"iconType"`
																	} `json:"icon"`
																	ServiceEndpoint struct {
																		ClickTrackingParams string `json:"clickTrackingParams"`
																		CommandMetadata     struct {
																			WebCommandMetadata struct {
																				SendPost bool `json:"sendPost"`
																			} `json:"webCommandMetadata"`
																		} `json:"commandMetadata"`
																		SignalServiceEndpoint struct {
																			Signal  string `json:"signal"`
																			Actions []struct {
																				ClickTrackingParams  string `json:"clickTrackingParams"`
																				AddToPlaylistCommand struct {
																					OpenMiniplayer      bool   `json:"openMiniplayer"`
																					VideoID             string `json:"videoId"`
																					ListType            string `json:"listType"`
																					OnCreateListCommand struct {
																						ClickTrackingParams string `json:"clickTrackingParams"`
																						CommandMetadata     struct {
																							WebCommandMetadata struct {
																								SendPost bool   `json:"sendPost"`
																								APIURL   string `json:"apiUrl"`
																							} `json:"webCommandMetadata"`
																						} `json:"commandMetadata"`
																						CreatePlaylistServiceEndpoint struct {
																							VideoIds []string `json:"videoIds"`
																							Params   string   `json:"params"`
																						} `json:"createPlaylistServiceEndpoint"`
																					} `json:"onCreateListCommand"`
																					VideoIds []string `json:"videoIds"`
																				} `json:"addToPlaylistCommand"`
																			} `json:"actions"`
																		} `json:"signalServiceEndpoint"`
																	} `json:"serviceEndpoint"`
																	TrackingParams string `json:"trackingParams"`
																} `json:"menuServiceItemRenderer,omitempty"`
															} `json:"items"`
															TrackingParams string `json:"trackingParams"`
															Accessibility  struct {
																AccessibilityData struct {
																	Label string `json:"label"`
																} `json:"accessibilityData"`
															} `json:"accessibility"`
														} `json:"menuRenderer"`
													} `json:"menu"`
													ThumbnailOverlays []struct {
														ThumbnailOverlayTimeStatusRenderer struct {
															Text struct {
																Accessibility struct {
																	AccessibilityData struct {
																		Label string `json:"label"`
																	} `json:"accessibilityData"`
																} `json:"accessibility"`
																SimpleText string `json:"simpleText"`
															} `json:"text"`
															Style string `json:"style"`
														} `json:"thumbnailOverlayTimeStatusRenderer,omitempty"`
														ThumbnailOverlayNowPlayingRenderer struct {
															Text struct {
																Runs []struct {
																	Text string `json:"text"`
																} `json:"runs"`
															} `json:"text"`
														} `json:"thumbnailOverlayNowPlayingRenderer,omitempty"`
													} `json:"thumbnailOverlays"`
													VideoInfo struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
													} `json:"videoInfo"`
												} `json:"playlistVideoRenderer,omitempty"`
												ContinuationItemRenderer struct {
													Trigger              string `json:"trigger"`
													ContinuationEndpoint struct {
														ClickTrackingParams string `json:"clickTrackingParams"`
														CommandMetadata     struct {
															WebCommandMetadata struct {
																SendPost bool   `json:"sendPost"`
																APIURL   string `json:"apiUrl"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
														ContinuationCommand struct {
															Token   string `json:"token"`
															Request string `json:"request"`
														} `json:"continuationCommand"`
													} `json:"continuationEndpoint"`
												} `json:"continuationItemRenderer,omitempty"`
											} `json:"contents"`
											PlaylistID     string `json:"playlistId"`
											IsEditable     bool   `json:"isEditable"`
											CanReorder     bool   `json:"canReorder"`
											TrackingParams string `json:"trackingParams"`
											TargetID       string `json:"targetId"`
										} `json:"playlistVideoListRenderer"`
									} `json:"contents"`
									TrackingParams string `json:"trackingParams"`
								} `json:"itemSectionRenderer"`
							} `json:"contents"`
							TrackingParams string `json:"trackingParams"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
					TrackingParams string `json:"trackingParams"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
	Header struct {
		PlaylistHeaderRenderer struct {
			PlaylistID string `json:"playlistId"`
			Title      struct {
				SimpleText string `json:"simpleText"`
			} `json:"title"`
			NumVideosText struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"numVideosText"`
			DescriptionText struct {
			} `json:"descriptionText"`
			OwnerText struct {
				Runs []struct {
					Text               string `json:"text"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
								APIURL      string `json:"apiUrl"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						BrowseEndpoint struct {
							BrowseID         string `json:"browseId"`
							CanonicalBaseURL string `json:"canonicalBaseUrl"`
						} `json:"browseEndpoint"`
					} `json:"navigationEndpoint"`
				} `json:"runs"`
			} `json:"ownerText"`
			ViewCountText struct {
				SimpleText string `json:"simpleText"`
			} `json:"viewCountText"`
			ShareData struct {
				CanShare bool `json:"canShare"`
			} `json:"shareData"`
			IsEditable    bool   `json:"isEditable"`
			Privacy       string `json:"privacy"`
			OwnerEndpoint struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				CommandMetadata     struct {
					WebCommandMetadata struct {
						URL         string `json:"url"`
						WebPageType string `json:"webPageType"`
						RootVe      int    `json:"rootVe"`
						APIURL      string `json:"apiUrl"`
					} `json:"webCommandMetadata"`
				} `json:"commandMetadata"`
				BrowseEndpoint struct {
					BrowseID         string `json:"browseId"`
					CanonicalBaseURL string `json:"canonicalBaseUrl"`
				} `json:"browseEndpoint"`
			} `json:"ownerEndpoint"`
			EditableDetails struct {
				CanDelete bool `json:"canDelete"`
			} `json:"editableDetails"`
			TrackingParams   string `json:"trackingParams"`
			ServiceEndpoints []struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				CommandMetadata     struct {
					WebCommandMetadata struct {
						SendPost bool   `json:"sendPost"`
						APIURL   string `json:"apiUrl"`
					} `json:"webCommandMetadata"`
				} `json:"commandMetadata"`
				PlaylistEditEndpoint struct {
					Actions []struct {
						Action           string `json:"action"`
						SourcePlaylistID string `json:"sourcePlaylistId"`
					} `json:"actions"`
				} `json:"playlistEditEndpoint"`
			} `json:"serviceEndpoints"`
			Stats []struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs,omitempty"`
				SimpleText string `json:"simpleText,omitempty"`
			} `json:"stats"`
			BriefStats []struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"briefStats"`
			PlaylistHeaderBanner struct {
				HeroPlaylistThumbnailRenderer struct {
					Thumbnail struct {
						Thumbnails []struct {
							URL    string `json:"url"`
							Width  int    `json:"width"`
							Height int    `json:"height"`
						} `json:"thumbnails"`
					} `json:"thumbnail"`
					MaxRatio       float64 `json:"maxRatio"`
					TrackingParams string  `json:"trackingParams"`
					OnTap          struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						WatchEndpoint struct {
							VideoID        string `json:"videoId"`
							PlaylistID     string `json:"playlistId"`
							PlayerParams   string `json:"playerParams"`
							LoggingContext struct {
								VssLoggingContext struct {
									SerializedContextData string `json:"serializedContextData"`
								} `json:"vssLoggingContext"`
							} `json:"loggingContext"`
							WatchEndpointSupportedOnesieConfig struct {
								HTML5PlaybackOnesieConfig struct {
									CommonConfig struct {
										URL string `json:"url"`
									} `json:"commonConfig"`
								} `json:"html5PlaybackOnesieConfig"`
							} `json:"watchEndpointSupportedOnesieConfig"`
						} `json:"watchEndpoint"`
					} `json:"onTap"`
					ThumbnailOverlays struct {
						ThumbnailOverlayHoverTextRenderer struct {
							Text struct {
								SimpleText string `json:"simpleText"`
							} `json:"text"`
							Icon struct {
								IconType string `json:"iconType"`
							} `json:"icon"`
						} `json:"thumbnailOverlayHoverTextRenderer"`
					} `json:"thumbnailOverlays"`
				} `json:"heroPlaylistThumbnailRenderer"`
			} `json:"playlistHeaderBanner"`
			SaveButton struct {
				ToggleButtonRenderer struct {
					Style struct {
						StyleType string `json:"styleType"`
					} `json:"style"`
					Size struct {
						SizeType string `json:"sizeType"`
					} `json:"size"`
					IsToggled   bool `json:"isToggled"`
					IsDisabled  bool `json:"isDisabled"`
					DefaultIcon struct {
						IconType string `json:"iconType"`
					} `json:"defaultIcon"`
					ToggledIcon struct {
						IconType string `json:"iconType"`
					} `json:"toggledIcon"`
					TrackingParams string `json:"trackingParams"`
					DefaultTooltip string `json:"defaultTooltip"`
					ToggledTooltip string `json:"toggledTooltip"`
					ToggledStyle   struct {
						StyleType string `json:"styleType"`
					} `json:"toggledStyle"`
					DefaultNavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								IgnoreNavigation bool `json:"ignoreNavigation"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						ModalEndpoint struct {
							Modal struct {
								ModalWithTitleAndButtonRenderer struct {
									Title struct {
										SimpleText string `json:"simpleText"`
									} `json:"title"`
									Content struct {
										SimpleText string `json:"simpleText"`
									} `json:"content"`
									Button struct {
										ButtonRenderer struct {
											Style      string `json:"style"`
											Size       string `json:"size"`
											IsDisabled bool   `json:"isDisabled"`
											Text       struct {
												SimpleText string `json:"simpleText"`
											} `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												CommandMetadata     struct {
													WebCommandMetadata struct {
														URL         string `json:"url"`
														WebPageType string `json:"webPageType"`
														RootVe      int    `json:"rootVe"`
													} `json:"webCommandMetadata"`
												} `json:"commandMetadata"`
												SignInEndpoint struct {
													NextEndpoint struct {
														ClickTrackingParams string `json:"clickTrackingParams"`
														CommandMetadata     struct {
															WebCommandMetadata struct {
																URL         string `json:"url"`
																WebPageType string `json:"webPageType"`
																RootVe      int    `json:"rootVe"`
																APIURL      string `json:"apiUrl"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
														BrowseEndpoint struct {
															BrowseID string `json:"browseId"`
														} `json:"browseEndpoint"`
													} `json:"nextEndpoint"`
													IdamTag string `json:"idamTag"`
												} `json:"signInEndpoint"`
											} `json:"navigationEndpoint"`
											TrackingParams string `json:"trackingParams"`
										} `json:"buttonRenderer"`
									} `json:"button"`
								} `json:"modalWithTitleAndButtonRenderer"`
							} `json:"modal"`
						} `json:"modalEndpoint"`
					} `json:"defaultNavigationEndpoint"`
					AccessibilityData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibilityData"`
					ToggledAccessibilityData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"toggledAccessibilityData"`
				} `json:"toggleButtonRenderer"`
			} `json:"saveButton"`
			ShareButton struct {
				ButtonRenderer struct {
					Style           string `json:"style"`
					Size            string `json:"size"`
					IsDisabled      bool   `json:"isDisabled"`
					ServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								SendPost bool   `json:"sendPost"`
								APIURL   string `json:"apiUrl"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						ShareEntityServiceEndpoint struct {
							SerializedShareEntity string `json:"serializedShareEntity"`
							Commands              []struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								OpenPopupAction     struct {
									Popup struct {
										UnifiedSharePanelRenderer struct {
											TrackingParams     string `json:"trackingParams"`
											ShowLoadingSpinner bool   `json:"showLoadingSpinner"`
										} `json:"unifiedSharePanelRenderer"`
									} `json:"popup"`
									PopupType string `json:"popupType"`
									BeReused  bool   `json:"beReused"`
								} `json:"openPopupAction"`
							} `json:"commands"`
						} `json:"shareEntityServiceEndpoint"`
					} `json:"serviceEndpoint"`
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					Tooltip           string `json:"tooltip"`
					TrackingParams    string `json:"trackingParams"`
					AccessibilityData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibilityData"`
				} `json:"buttonRenderer"`
			} `json:"shareButton"`
			MoreActionsMenu struct {
				MenuRenderer struct {
					Items []struct {
						MenuNavigationItemRenderer struct {
							Text struct {
								SimpleText string `json:"simpleText"`
							} `json:"text"`
							Icon struct {
								IconType string `json:"iconType"`
							} `json:"icon"`
							NavigationEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										URL         string `json:"url"`
										WebPageType string `json:"webPageType"`
										RootVe      int    `json:"rootVe"`
										APIURL      string `json:"apiUrl"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								BrowseEndpoint struct {
									BrowseID       string `json:"browseId"`
									Params         string `json:"params"`
									Nofollow       bool   `json:"nofollow"`
									NavigationType string `json:"navigationType"`
								} `json:"browseEndpoint"`
							} `json:"navigationEndpoint"`
							TrackingParams string `json:"trackingParams"`
						} `json:"menuNavigationItemRenderer"`
					} `json:"items"`
					TrackingParams string `json:"trackingParams"`
					Accessibility  struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibility"`
					TargetID string `json:"targetId"`
				} `json:"menuRenderer"`
			} `json:"moreActionsMenu"`
			PlayButton struct {
				ButtonRenderer struct {
					Style      string `json:"style"`
					Size       string `json:"size"`
					IsDisabled bool   `json:"isDisabled"`
					Text       struct {
						SimpleText string `json:"simpleText"`
					} `json:"text"`
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						WatchEndpoint struct {
							VideoID        string `json:"videoId"`
							PlaylistID     string `json:"playlistId"`
							PlayerParams   string `json:"playerParams"`
							LoggingContext struct {
								VssLoggingContext struct {
									SerializedContextData string `json:"serializedContextData"`
								} `json:"vssLoggingContext"`
							} `json:"loggingContext"`
							WatchEndpointSupportedOnesieConfig struct {
								HTML5PlaybackOnesieConfig struct {
									CommonConfig struct {
										URL string `json:"url"`
									} `json:"commonConfig"`
								} `json:"html5PlaybackOnesieConfig"`
							} `json:"watchEndpointSupportedOnesieConfig"`
						} `json:"watchEndpoint"`
					} `json:"navigationEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"buttonRenderer"`
			} `json:"playButton"`
			ShufflePlayButton struct {
				ButtonRenderer struct {
					Style      string `json:"style"`
					Size       string `json:"size"`
					IsDisabled bool   `json:"isDisabled"`
					Text       struct {
						SimpleText string `json:"simpleText"`
					} `json:"text"`
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						WatchEndpoint struct {
							VideoID        string `json:"videoId"`
							PlaylistID     string `json:"playlistId"`
							Params         string `json:"params"`
							PlayerParams   string `json:"playerParams"`
							LoggingContext struct {
								VssLoggingContext struct {
									SerializedContextData string `json:"serializedContextData"`
								} `json:"vssLoggingContext"`
							} `json:"loggingContext"`
							WatchEndpointSupportedOnesieConfig struct {
								HTML5PlaybackOnesieConfig struct {
									CommonConfig struct {
										URL string `json:"url"`
									} `json:"commonConfig"`
								} `json:"html5PlaybackOnesieConfig"`
							} `json:"watchEndpointSupportedOnesieConfig"`
						} `json:"watchEndpoint"`
					} `json:"navigationEndpoint"`
					TrackingParams string `json:"trackingParams"`
				} `json:"buttonRenderer"`
			} `json:"shufflePlayButton"`
			OnDescriptionTap struct {
				ClickTrackingParams string `json:"clickTrackingParams"`
				OpenPopupAction     struct {
					Popup struct {
						FancyDismissibleDialogRenderer struct {
							DialogMessage struct {
							} `json:"dialogMessage"`
							Title struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"title"`
							ConfirmLabel struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"confirmLabel"`
							TrackingParams string `json:"trackingParams"`
						} `json:"fancyDismissibleDialogRenderer"`
					} `json:"popup"`
					PopupType string `json:"popupType"`
				} `json:"openPopupAction"`
			} `json:"onDescriptionTap"`
			CinematicContainer struct {
				CinematicContainerRenderer struct {
					BackgroundImageConfig struct {
						Thumbnail struct {
							Thumbnails []struct {
								URL    string `json:"url"`
								Width  int    `json:"width"`
								Height int    `json:"height"`
							} `json:"thumbnails"`
						} `json:"thumbnail"`
					} `json:"backgroundImageConfig"`
					GradientColorConfig []struct {
						LightThemeColor int64 `json:"lightThemeColor"`
						DarkThemeColor  int64 `json:"darkThemeColor"`
						StartLocation   int   `json:"startLocation"`
					} `json:"gradientColorConfig"`
					Config struct {
						LightThemeBackgroundColor int64 `json:"lightThemeBackgroundColor"`
						DarkThemeBackgroundColor  int64 `json:"darkThemeBackgroundColor"`
						ColorSourceSizeMultiplier int   `json:"colorSourceSizeMultiplier"`
						ApplyClientImageBlur      bool  `json:"applyClientImageBlur"`
					} `json:"config"`
				} `json:"cinematicContainerRenderer"`
			} `json:"cinematicContainer"`
			Byline []struct {
				PlaylistBylineRenderer struct {
					Text struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
				} `json:"playlistBylineRenderer"`
			} `json:"byline"`
			DescriptionTapText struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"descriptionTapText"`
		} `json:"playlistHeaderRenderer"`
	} `json:"header"`
	Alerts []struct {
		AlertWithButtonRenderer struct {
			Type string `json:"type"`
			Text struct {
				SimpleText string `json:"simpleText"`
			} `json:"text"`
			DismissButton struct {
				ButtonRenderer struct {
					Style      string `json:"style"`
					Size       string `json:"size"`
					IsDisabled bool   `json:"isDisabled"`
					Icon       struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					TrackingParams    string `json:"trackingParams"`
					AccessibilityData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibilityData"`
				} `json:"buttonRenderer"`
			} `json:"dismissButton"`
		} `json:"alertWithButtonRenderer"`
	} `json:"alerts"`
	Metadata struct {
		PlaylistMetadataRenderer struct {
			Title                  string `json:"title"`
			AndroidAppindexingLink string `json:"androidAppindexingLink"`
			IosAppindexingLink     string `json:"iosAppindexingLink"`
		} `json:"playlistMetadataRenderer"`
	} `json:"metadata"`
	TrackingParams string `json:"trackingParams"`
	Topbar         struct {
		DesktopTopbarRenderer struct {
			Logo struct {
				TopbarLogoRenderer struct {
					IconImage struct {
						IconType string `json:"iconType"`
					} `json:"iconImage"`
					TooltipText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"tooltipText"`
					Endpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
								APIURL      string `json:"apiUrl"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						BrowseEndpoint struct {
							BrowseID string `json:"browseId"`
						} `json:"browseEndpoint"`
					} `json:"endpoint"`
					TrackingParams    string `json:"trackingParams"`
					OverrideEntityKey string `json:"overrideEntityKey"`
				} `json:"topbarLogoRenderer"`
			} `json:"logo"`
			Searchbox struct {
				FusionSearchboxRenderer struct {
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					PlaceholderText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"placeholderText"`
					Config struct {
						WebSearchboxConfig struct {
							RequestLanguage     string `json:"requestLanguage"`
							RequestDomain       string `json:"requestDomain"`
							HasOnscreenKeyboard bool   `json:"hasOnscreenKeyboard"`
							FocusSearchbox      bool   `json:"focusSearchbox"`
						} `json:"webSearchboxConfig"`
					} `json:"config"`
					TrackingParams string `json:"trackingParams"`
					SearchEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						SearchEndpoint struct {
							Query string `json:"query"`
						} `json:"searchEndpoint"`
					} `json:"searchEndpoint"`
					ClearButton struct {
						ButtonRenderer struct {
							Style      string `json:"style"`
							Size       string `json:"size"`
							IsDisabled bool   `json:"isDisabled"`
							Icon       struct {
								IconType string `json:"iconType"`
							} `json:"icon"`
							TrackingParams    string `json:"trackingParams"`
							AccessibilityData struct {
								AccessibilityData struct {
									Label string `json:"label"`
								} `json:"accessibilityData"`
							} `json:"accessibilityData"`
						} `json:"buttonRenderer"`
					} `json:"clearButton"`
				} `json:"fusionSearchboxRenderer"`
			} `json:"searchbox"`
			TrackingParams string `json:"trackingParams"`
			CountryCode    string `json:"countryCode"`
			TopbarButtons  []struct {
				TopbarMenuButtonRenderer struct {
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					MenuRequest struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								SendPost bool   `json:"sendPost"`
								APIURL   string `json:"apiUrl"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						SignalServiceEndpoint struct {
							Signal  string `json:"signal"`
							Actions []struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								OpenPopupAction     struct {
									Popup struct {
										MultiPageMenuRenderer struct {
											TrackingParams     string `json:"trackingParams"`
											Style              string `json:"style"`
											ShowLoadingSpinner bool   `json:"showLoadingSpinner"`
										} `json:"multiPageMenuRenderer"`
									} `json:"popup"`
									PopupType string `json:"popupType"`
									BeReused  bool   `json:"beReused"`
								} `json:"openPopupAction"`
							} `json:"actions"`
						} `json:"signalServiceEndpoint"`
					} `json:"menuRequest"`
					TrackingParams string `json:"trackingParams"`
					Accessibility  struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibility"`
					Tooltip string `json:"tooltip"`
					Style   string `json:"style"`
				} `json:"topbarMenuButtonRenderer,omitempty"`
				ButtonRenderer struct {
					Style string `json:"style"`
					Size  string `json:"size"`
					Text  struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						SignInEndpoint struct {
							IdamTag string `json:"idamTag"`
						} `json:"signInEndpoint"`
					} `json:"navigationEndpoint"`
					TrackingParams string `json:"trackingParams"`
					TargetID       string `json:"targetId"`
				} `json:"buttonRenderer,omitempty"`
			} `json:"topbarButtons"`
			HotkeyDialog struct {
				HotkeyDialogRenderer struct {
					Title struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"title"`
					Sections []struct {
						HotkeyDialogSectionRenderer struct {
							Title struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"title"`
							Options []struct {
								HotkeyDialogSectionOptionRenderer struct {
									Label struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"label"`
									Hotkey string `json:"hotkey"`
								} `json:"hotkeyDialogSectionOptionRenderer,omitempty"`
							} `json:"options"`
						} `json:"hotkeyDialogSectionRenderer"`
					} `json:"sections"`
					DismissButton struct {
						ButtonRenderer struct {
							Style      string `json:"style"`
							Size       string `json:"size"`
							IsDisabled bool   `json:"isDisabled"`
							Text       struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"text"`
							TrackingParams string `json:"trackingParams"`
						} `json:"buttonRenderer"`
					} `json:"dismissButton"`
					TrackingParams string `json:"trackingParams"`
				} `json:"hotkeyDialogRenderer"`
			} `json:"hotkeyDialog"`
			BackButton struct {
				ButtonRenderer struct {
					TrackingParams string `json:"trackingParams"`
					Command        struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								SendPost bool `json:"sendPost"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						SignalServiceEndpoint struct {
							Signal  string `json:"signal"`
							Actions []struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								SignalAction        struct {
									Signal string `json:"signal"`
								} `json:"signalAction"`
							} `json:"actions"`
						} `json:"signalServiceEndpoint"`
					} `json:"command"`
				} `json:"buttonRenderer"`
			} `json:"backButton"`
			ForwardButton struct {
				ButtonRenderer struct {
					TrackingParams string `json:"trackingParams"`
					Command        struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								SendPost bool `json:"sendPost"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						SignalServiceEndpoint struct {
							Signal  string `json:"signal"`
							Actions []struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								SignalAction        struct {
									Signal string `json:"signal"`
								} `json:"signalAction"`
							} `json:"actions"`
						} `json:"signalServiceEndpoint"`
					} `json:"command"`
				} `json:"buttonRenderer"`
			} `json:"forwardButton"`
			A11YSkipNavigationButton struct {
				ButtonRenderer struct {
					Style      string `json:"style"`
					Size       string `json:"size"`
					IsDisabled bool   `json:"isDisabled"`
					Text       struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"text"`
					TrackingParams string `json:"trackingParams"`
					Command        struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								SendPost bool `json:"sendPost"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						SignalServiceEndpoint struct {
							Signal  string `json:"signal"`
							Actions []struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								SignalAction        struct {
									Signal string `json:"signal"`
								} `json:"signalAction"`
							} `json:"actions"`
						} `json:"signalServiceEndpoint"`
					} `json:"command"`
				} `json:"buttonRenderer"`
			} `json:"a11ySkipNavigationButton"`
			VoiceSearchButton struct {
				ButtonRenderer struct {
					Style           string `json:"style"`
					Size            string `json:"size"`
					IsDisabled      bool   `json:"isDisabled"`
					ServiceEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								SendPost bool `json:"sendPost"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						SignalServiceEndpoint struct {
							Signal  string `json:"signal"`
							Actions []struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								OpenPopupAction     struct {
									Popup struct {
										VoiceSearchDialogRenderer struct {
											PlaceholderHeader struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"placeholderHeader"`
											PromptHeader struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"promptHeader"`
											ExampleQuery1 struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"exampleQuery1"`
											ExampleQuery2 struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"exampleQuery2"`
											PromptMicrophoneLabel struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"promptMicrophoneLabel"`
											LoadingHeader struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"loadingHeader"`
											ConnectionErrorHeader struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"connectionErrorHeader"`
											ConnectionErrorMicrophoneLabel struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"connectionErrorMicrophoneLabel"`
											PermissionsHeader struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"permissionsHeader"`
											PermissionsSubtext struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"permissionsSubtext"`
											DisabledHeader struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"disabledHeader"`
											DisabledSubtext struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"disabledSubtext"`
											MicrophoneButtonAriaLabel struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"microphoneButtonAriaLabel"`
											ExitButton struct {
												ButtonRenderer struct {
													Style      string `json:"style"`
													Size       string `json:"size"`
													IsDisabled bool   `json:"isDisabled"`
													Icon       struct {
														IconType string `json:"iconType"`
													} `json:"icon"`
													TrackingParams    string `json:"trackingParams"`
													AccessibilityData struct {
														AccessibilityData struct {
															Label string `json:"label"`
														} `json:"accessibilityData"`
													} `json:"accessibilityData"`
												} `json:"buttonRenderer"`
											} `json:"exitButton"`
											TrackingParams            string `json:"trackingParams"`
											MicrophoneOffPromptHeader struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"microphoneOffPromptHeader"`
										} `json:"voiceSearchDialogRenderer"`
									} `json:"popup"`
									PopupType string `json:"popupType"`
								} `json:"openPopupAction"`
							} `json:"actions"`
						} `json:"signalServiceEndpoint"`
					} `json:"serviceEndpoint"`
					Icon struct {
						IconType string `json:"iconType"`
					} `json:"icon"`
					Tooltip           string `json:"tooltip"`
					TrackingParams    string `json:"trackingParams"`
					AccessibilityData struct {
						AccessibilityData struct {
							Label string `json:"label"`
						} `json:"accessibilityData"`
					} `json:"accessibilityData"`
				} `json:"buttonRenderer"`
			} `json:"voiceSearchButton"`
		} `json:"desktopTopbarRenderer"`
	} `json:"topbar"`
	Microformat struct {
		MicroformatDataRenderer struct {
			URLCanonical string `json:"urlCanonical"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			Thumbnail    struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			SiteName           string `json:"siteName"`
			AppName            string `json:"appName"`
			AndroidPackage     string `json:"androidPackage"`
			IosAppStoreID      string `json:"iosAppStoreId"`
			IosAppArguments    string `json:"iosAppArguments"`
			OgType             string `json:"ogType"`
			URLApplinksWeb     string `json:"urlApplinksWeb"`
			URLApplinksIos     string `json:"urlApplinksIos"`
			URLApplinksAndroid string `json:"urlApplinksAndroid"`
			URLTwitterIos      string `json:"urlTwitterIos"`
			URLTwitterAndroid  string `json:"urlTwitterAndroid"`
			TwitterCardType    string `json:"twitterCardType"`
			TwitterSiteHandle  string `json:"twitterSiteHandle"`
			SchemaDotOrgType   string `json:"schemaDotOrgType"`
			Noindex            bool   `json:"noindex"`
			Unlisted           bool   `json:"unlisted"`
			LinkAlternates     []struct {
				HrefURL string `json:"hrefUrl"`
			} `json:"linkAlternates"`
		} `json:"microformatDataRenderer"`
	} `json:"microformat"`
	Sidebar struct {
		PlaylistSidebarRenderer struct {
			Items []struct {
				PlaylistSidebarPrimaryInfoRenderer struct {
					ThumbnailRenderer struct {
						PlaylistVideoThumbnailRenderer struct {
							Thumbnail struct {
								Thumbnails []struct {
									URL    string `json:"url"`
									Width  int    `json:"width"`
									Height int    `json:"height"`
								} `json:"thumbnails"`
							} `json:"thumbnail"`
							TrackingParams string `json:"trackingParams"`
						} `json:"playlistVideoThumbnailRenderer"`
					} `json:"thumbnailRenderer"`
					Title struct {
						Runs []struct {
							Text               string `json:"text"`
							NavigationEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										URL         string `json:"url"`
										WebPageType string `json:"webPageType"`
										RootVe      int    `json:"rootVe"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								WatchEndpoint struct {
									VideoID        string `json:"videoId"`
									PlaylistID     string `json:"playlistId"`
									PlayerParams   string `json:"playerParams"`
									LoggingContext struct {
										VssLoggingContext struct {
											SerializedContextData string `json:"serializedContextData"`
										} `json:"vssLoggingContext"`
									} `json:"loggingContext"`
									WatchEndpointSupportedOnesieConfig struct {
										HTML5PlaybackOnesieConfig struct {
											CommonConfig struct {
												URL string `json:"url"`
											} `json:"commonConfig"`
										} `json:"html5PlaybackOnesieConfig"`
									} `json:"watchEndpointSupportedOnesieConfig"`
								} `json:"watchEndpoint"`
							} `json:"navigationEndpoint"`
						} `json:"runs"`
					} `json:"title"`
					Stats []struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs,omitempty"`
						SimpleText string `json:"simpleText,omitempty"`
					} `json:"stats"`
					Menu struct {
						MenuRenderer struct {
							Items []struct {
								MenuNavigationItemRenderer struct {
									Text struct {
										SimpleText string `json:"simpleText"`
									} `json:"text"`
									Icon struct {
										IconType string `json:"iconType"`
									} `json:"icon"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												URL         string `json:"url"`
												WebPageType string `json:"webPageType"`
												RootVe      int    `json:"rootVe"`
												APIURL      string `json:"apiUrl"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										BrowseEndpoint struct {
											BrowseID       string `json:"browseId"`
											Params         string `json:"params"`
											Nofollow       bool   `json:"nofollow"`
											NavigationType string `json:"navigationType"`
										} `json:"browseEndpoint"`
									} `json:"navigationEndpoint"`
									TrackingParams string `json:"trackingParams"`
								} `json:"menuNavigationItemRenderer"`
							} `json:"items"`
							TrackingParams  string `json:"trackingParams"`
							TopLevelButtons []struct {
								ToggleButtonRenderer struct {
									Style struct {
										StyleType string `json:"styleType"`
									} `json:"style"`
									Size struct {
										SizeType string `json:"sizeType"`
									} `json:"size"`
									IsToggled   bool `json:"isToggled"`
									IsDisabled  bool `json:"isDisabled"`
									DefaultIcon struct {
										IconType string `json:"iconType"`
									} `json:"defaultIcon"`
									ToggledIcon struct {
										IconType string `json:"iconType"`
									} `json:"toggledIcon"`
									TrackingParams            string `json:"trackingParams"`
									DefaultTooltip            string `json:"defaultTooltip"`
									ToggledTooltip            string `json:"toggledTooltip"`
									DefaultNavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												IgnoreNavigation bool `json:"ignoreNavigation"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										ModalEndpoint struct {
											Modal struct {
												ModalWithTitleAndButtonRenderer struct {
													Title struct {
														SimpleText string `json:"simpleText"`
													} `json:"title"`
													Content struct {
														SimpleText string `json:"simpleText"`
													} `json:"content"`
													Button struct {
														ButtonRenderer struct {
															Style      string `json:"style"`
															Size       string `json:"size"`
															IsDisabled bool   `json:"isDisabled"`
															Text       struct {
																SimpleText string `json:"simpleText"`
															} `json:"text"`
															NavigationEndpoint struct {
																ClickTrackingParams string `json:"clickTrackingParams"`
																CommandMetadata     struct {
																	WebCommandMetadata struct {
																		URL         string `json:"url"`
																		WebPageType string `json:"webPageType"`
																		RootVe      int    `json:"rootVe"`
																	} `json:"webCommandMetadata"`
																} `json:"commandMetadata"`
																SignInEndpoint struct {
																	NextEndpoint struct {
																		ClickTrackingParams string `json:"clickTrackingParams"`
																		CommandMetadata     struct {
																			WebCommandMetadata struct {
																				URL         string `json:"url"`
																				WebPageType string `json:"webPageType"`
																				RootVe      int    `json:"rootVe"`
																				APIURL      string `json:"apiUrl"`
																			} `json:"webCommandMetadata"`
																		} `json:"commandMetadata"`
																		BrowseEndpoint struct {
																			BrowseID string `json:"browseId"`
																		} `json:"browseEndpoint"`
																	} `json:"nextEndpoint"`
																	IdamTag string `json:"idamTag"`
																} `json:"signInEndpoint"`
															} `json:"navigationEndpoint"`
															TrackingParams string `json:"trackingParams"`
														} `json:"buttonRenderer"`
													} `json:"button"`
												} `json:"modalWithTitleAndButtonRenderer"`
											} `json:"modal"`
										} `json:"modalEndpoint"`
									} `json:"defaultNavigationEndpoint"`
									AccessibilityData struct {
										AccessibilityData struct {
											Label string `json:"label"`
										} `json:"accessibilityData"`
									} `json:"accessibilityData"`
									ToggledAccessibilityData struct {
										AccessibilityData struct {
											Label string `json:"label"`
										} `json:"accessibilityData"`
									} `json:"toggledAccessibilityData"`
								} `json:"toggleButtonRenderer,omitempty"`
								ButtonRenderer struct {
									Style      string `json:"style"`
									Size       string `json:"size"`
									IsDisabled bool   `json:"isDisabled"`
									Icon       struct {
										IconType string `json:"iconType"`
									} `json:"icon"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												URL         string `json:"url"`
												WebPageType string `json:"webPageType"`
												RootVe      int    `json:"rootVe"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										WatchEndpoint struct {
											VideoID        string `json:"videoId"`
											PlaylistID     string `json:"playlistId"`
											Params         string `json:"params"`
											PlayerParams   string `json:"playerParams"`
											LoggingContext struct {
												VssLoggingContext struct {
													SerializedContextData string `json:"serializedContextData"`
												} `json:"vssLoggingContext"`
											} `json:"loggingContext"`
											WatchEndpointSupportedOnesieConfig struct {
												HTML5PlaybackOnesieConfig struct {
													CommonConfig struct {
														URL string `json:"url"`
													} `json:"commonConfig"`
												} `json:"html5PlaybackOnesieConfig"`
											} `json:"watchEndpointSupportedOnesieConfig"`
										} `json:"watchEndpoint"`
									} `json:"navigationEndpoint"`
									Accessibility struct {
										Label string `json:"label"`
									} `json:"accessibility"`
									Tooltip        string `json:"tooltip"`
									TrackingParams string `json:"trackingParams"`
								} `json:"buttonRenderer,omitempty"`
							} `json:"topLevelButtons"`
							Accessibility struct {
								AccessibilityData struct {
									Label string `json:"label"`
								} `json:"accessibilityData"`
							} `json:"accessibility"`
							TargetID string `json:"targetId"`
						} `json:"menuRenderer"`
					} `json:"menu"`
					ThumbnailOverlays []struct {
						ThumbnailOverlaySidePanelRenderer struct {
							Text struct {
								SimpleText string `json:"simpleText"`
							} `json:"text"`
							Icon struct {
								IconType string `json:"iconType"`
							} `json:"icon"`
						} `json:"thumbnailOverlaySidePanelRenderer"`
					} `json:"thumbnailOverlays"`
					NavigationEndpoint struct {
						ClickTrackingParams string `json:"clickTrackingParams"`
						CommandMetadata     struct {
							WebCommandMetadata struct {
								URL         string `json:"url"`
								WebPageType string `json:"webPageType"`
								RootVe      int    `json:"rootVe"`
							} `json:"webCommandMetadata"`
						} `json:"commandMetadata"`
						WatchEndpoint struct {
							VideoID        string `json:"videoId"`
							PlaylistID     string `json:"playlistId"`
							PlayerParams   string `json:"playerParams"`
							LoggingContext struct {
								VssLoggingContext struct {
									SerializedContextData string `json:"serializedContextData"`
								} `json:"vssLoggingContext"`
							} `json:"loggingContext"`
							WatchEndpointSupportedOnesieConfig struct {
								HTML5PlaybackOnesieConfig struct {
									CommonConfig struct {
										URL string `json:"url"`
									} `json:"commonConfig"`
								} `json:"html5PlaybackOnesieConfig"`
							} `json:"watchEndpointSupportedOnesieConfig"`
						} `json:"watchEndpoint"`
					} `json:"navigationEndpoint"`
					Description struct {
					} `json:"description"`
					ShowMoreText struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"showMoreText"`
				} `json:"playlistSidebarPrimaryInfoRenderer,omitempty"`
				PlaylistSidebarSecondaryInfoRenderer struct {
					VideoOwner struct {
						VideoOwnerRenderer struct {
							Thumbnail struct {
								Thumbnails []struct {
									URL    string `json:"url"`
									Width  int    `json:"width"`
									Height int    `json:"height"`
								} `json:"thumbnails"`
							} `json:"thumbnail"`
							Title struct {
								Runs []struct {
									Text               string `json:"text"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												URL         string `json:"url"`
												WebPageType string `json:"webPageType"`
												RootVe      int    `json:"rootVe"`
												APIURL      string `json:"apiUrl"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										BrowseEndpoint struct {
											BrowseID         string `json:"browseId"`
											CanonicalBaseURL string `json:"canonicalBaseUrl"`
										} `json:"browseEndpoint"`
									} `json:"navigationEndpoint"`
								} `json:"runs"`
							} `json:"title"`
							NavigationEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										URL         string `json:"url"`
										WebPageType string `json:"webPageType"`
										RootVe      int    `json:"rootVe"`
										APIURL      string `json:"apiUrl"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								BrowseEndpoint struct {
									BrowseID         string `json:"browseId"`
									CanonicalBaseURL string `json:"canonicalBaseUrl"`
								} `json:"browseEndpoint"`
							} `json:"navigationEndpoint"`
							TrackingParams string `json:"trackingParams"`
						} `json:"videoOwnerRenderer"`
					} `json:"videoOwner"`
					Button struct {
						ButtonRenderer struct {
							Style      string `json:"style"`
							Size       string `json:"size"`
							IsDisabled bool   `json:"isDisabled"`
							Text       struct {
								Runs []struct {
									Text string `json:"text"`
								} `json:"runs"`
							} `json:"text"`
							NavigationEndpoint struct {
								ClickTrackingParams string `json:"clickTrackingParams"`
								CommandMetadata     struct {
									WebCommandMetadata struct {
										IgnoreNavigation bool `json:"ignoreNavigation"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
								ModalEndpoint struct {
									Modal struct {
										ModalWithTitleAndButtonRenderer struct {
											Title struct {
												SimpleText string `json:"simpleText"`
											} `json:"title"`
											Content struct {
												SimpleText string `json:"simpleText"`
											} `json:"content"`
											Button struct {
												ButtonRenderer struct {
													Style      string `json:"style"`
													Size       string `json:"size"`
													IsDisabled bool   `json:"isDisabled"`
													Text       struct {
														SimpleText string `json:"simpleText"`
													} `json:"text"`
													NavigationEndpoint struct {
														ClickTrackingParams string `json:"clickTrackingParams"`
														CommandMetadata     struct {
															WebCommandMetadata struct {
																URL         string `json:"url"`
																WebPageType string `json:"webPageType"`
																RootVe      int    `json:"rootVe"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
														SignInEndpoint struct {
															NextEndpoint struct {
																ClickTrackingParams string `json:"clickTrackingParams"`
																CommandMetadata     struct {
																	WebCommandMetadata struct {
																		URL         string `json:"url"`
																		WebPageType string `json:"webPageType"`
																		RootVe      int    `json:"rootVe"`
																		APIURL      string `json:"apiUrl"`
																	} `json:"webCommandMetadata"`
																} `json:"commandMetadata"`
																BrowseEndpoint struct {
																	BrowseID string `json:"browseId"`
																} `json:"browseEndpoint"`
															} `json:"nextEndpoint"`
															ContinueAction string `json:"continueAction"`
															IdamTag        string `json:"idamTag"`
														} `json:"signInEndpoint"`
													} `json:"navigationEndpoint"`
													TrackingParams string `json:"trackingParams"`
												} `json:"buttonRenderer"`
											} `json:"button"`
										} `json:"modalWithTitleAndButtonRenderer"`
									} `json:"modal"`
								} `json:"modalEndpoint"`
							} `json:"navigationEndpoint"`
							TrackingParams string `json:"trackingParams"`
						} `json:"buttonRenderer"`
					} `json:"button"`
				} `json:"playlistSidebarSecondaryInfoRenderer,omitempty"`
			} `json:"items"`
			TrackingParams string `json:"trackingParams"`
		} `json:"playlistSidebarRenderer"`
	} `json:"sidebar"`
}
