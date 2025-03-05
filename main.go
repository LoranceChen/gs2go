package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	patch_data "gs2go/ck-patch-data-protocol/protobuf"
	patch_data2 "gs2go/ck-patch-data-protocol2/protobuf"
	"gs2go/proto_define"
	"gs2go/router"
	"net/http"
	"os"
	"text/template"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"google.golang.org/protobuf/proto"

	protojson "google.golang.org/protobuf/encoding/protojson"

	"github.com/gorilla/websocket"
	"github.com/grafana/pyroscope-go"
	"github.com/rs/zerolog/log"
)

var addr = flag.String("addr", "localhost:8880", "http service address")

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info().Msgf("upgrade: ", err)
		return
	}

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Info().Msgf("readMessage: %v", err)
		}

		log.Info().Msgf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Info().Msgf("writeMessage: %v", err)
			break
		}
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	testPatchData()

	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "simple.app",

		// replace this with the address of pyroscope server
		ServerAddress: "http://kmc.aow.com",

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// Optional HTTP Basic authentication (Grafana Cloud)
		// BasicAuthUser:     "<USER>",
		BasicAuthPassword: "TODO",
		// Optional Pyroscope tenant ID (only needed if using multi-tenancy). Not needed for Grafana Cloud.
		// TenantID:          "<TenantID>",

		// by default all profilers are enabled,
		// but you can select the ones you want to use:
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})

	if err != nil {
		log.Error().Msgf("profiler start fail: ", err)
	}

	//a1 := a.A1{DataA: 1}
	//
	//a1.GetB1()

	ctx := context.Background()
	// createDB
	//pg := createDB()
	_, err = pgxInit(ctx, "TODO")

	// proto test
	request := proto_define.HelloRequest{
		Msg:      "go proto",
		Sequence: 0,
	}
	marshal, err := protojson.Marshal(&request)
	if err != nil {
		log.Info().Msgf("fail marshal pb: ", err)
		return
	}
	log.Info().Msgf("hello request: %v", string(marshal))

	response1 := proto_define.SignUpResponse{
		Kingdom: &proto_define.Kingdom{
			Id:    0,
			Name:  "",
			Items: nil,
		}, Name: "123"}

	_, _ = proto.Marshal(&response1)

	marshal2, err := protojson.Marshal(&response1)
	if err != nil {
		log.Printf("fail marshal pb: ", err)
		return
	}

	log.Info().Msgf("SignUpResponse: %v", string(marshal2))
	marshal3, err := json.Marshal(response1)
	log.Info().Msgf("SignUpResponse simple json : %v", string(marshal3))

	log.Info().Msgf("starting websocket & protobuf server on %s", *addr)
	flag.Parse()
	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)

	http.HandleFunc("/wspb", wspb)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		panic(err)
	}
	// log.Fatal(http.ListenAndServe(*addr, nil))

}

func createDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
	return conn
}

func dbPool() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}

func wspb(w http.ResponseWriter, r *http.Request) {

	_ = patch_data.NullableInt{Value: &patch_data.NullableInt_Int{Int: 10}}
	router.WsPbRouter(w, r, upgrader)
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))

func testPatchData() {
	//_ = patch_data.Date{
	//	Year:  0,
	//	Month: 0,
	//	Day:   0,
	//}

	rst := patch_data.ZonePaperShop{
		PaperItemId:    0,
		PaperId:        "",
		PaperName:      "",
		PaperThumb:     "",
		PaperSize:      0,
		Author:         "",
		ThemeSeasonId:  0,
		PriceItemId:    nil,
		PriceQuantity:  nil,
		PriceItemId2:   nil,
		PriceQuantity2: nil,
		Removed:        false,
	}

	adr := func(s string) *string { return &s }

	tstring := patch_data2.TString_builder{
		Key:          adr(""),
		OriginalText: adr(""),
	}.Build()

	myId := int32(1)

	description := patch_data2.ShopDescription_builder{
		Id:          &myId,
		Description: tstring,
	}.Build()

	a, e := proto.Marshal(description)
	if e != nil {
		println("aaaa: %s", a)
	}
	fmt.Println(rst)
	fmt.Println(description)

}

func pgxInit(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	conn, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	println("asd")
	return conn, nil
}
