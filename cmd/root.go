package cmd

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"fmt"
	"os"
	"github.com/gorilla/mux"
	"net/http"
	ws "github.com/renosyah/basicWebSocket/ws"
)

var (
	cfgFile string
	isNotListening = true
)
var rootCmd = &cobra.Command{
	Use: "basic",
	Run: func(cmd *cobra.Command, args []string) {
		r := mux.NewRouter()
		http.Handle("/", r)

		h := ws.NewHub()
		r.Handle("/ws", ws.WsHandler{Hub : h})
		r.HandleFunc("/ping",ws.StartSendMessagesAsClient)
		r.HandleFunc("/start",func(http.ResponseWriter,*http.Request){
			if isNotListening {
				fmt.Println("listenig as websocket client...")
				go ws.StartListeningMessagesAsClient()
			}
			isNotListening = false
		})

		if err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("app.port")), r); err != nil {
			fmt.Println("failed to serve: ", err.Error())
		}
	},
}
func init() {
	cobra.OnInitialize(initConfig)
}

func Execute(){
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {

	viper.SetConfigType("toml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)

	} else {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/basicWebSocket")
		viper.SetConfigName(".basicWebSocket")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}