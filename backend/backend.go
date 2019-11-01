package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"rpcx-demo/service/product/model"
	modeluser "rpcx-demo/service/user/model"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/smallnest/rpcx/client"
)

var (
	addr  = flag.String("addr", ":8080", "http address")
	paddr = flag.String("product-image-addr", "localhost:8972", "图片服务地址")
)

var (
	xclient    client.XClient
	userclinet client.XClient
)

func main() {
	d := client.NewPeer2PeerDiscovery("tcp@"+*paddr, "")
	xclient = client.NewXClient("ProductImage", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
	userclinet = client.NewXClient("auth", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer userclinet.Close()

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/banner-ppl.png", index)
	router.GET("/banner-ppl-women.png", index)
	router.GET("/bk-sale.png", index)
	router.ServeFiles("/_nuxt/*filepath", http.Dir("../web/dist/_nuxt"))
	router.ServeFiles("/cart/*filepath", http.Dir("../web/dist/cart"))
	router.ServeFiles("/men/*filepath", http.Dir("../web/dist/men"))
	router.ServeFiles("/sale/*filepath", http.Dir("../web/dist/sale"))
	router.ServeFiles("/women/*filepath", http.Dir("../web/dist/women"))

	// product
	router.GET("/products_images/:name", productsImages)

	//auth
	router.POST("/auth", auth)

	//user
	router.GET("/user/:say", say)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func say(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	say := ps.ByName("say")
	resp := new(modeluser.SayResponse)
	req := modeluser.SayRequest(say)
	err := userclinet.Call(context.Background(), "Say", req, resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("*resp:", *resp)
	resp_byte, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp_byte)
}

func auth(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	pwd := ps.ByName("password")
	resp := &modeluser.AuthResponse{}
	req := modeluser.AuthRequest{
		UserName: name,
		Password: pwd,
	}
	err := userclinet.Call(context.Background(), "Login", req, resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("*resp:", *resp)
	resp_byte, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp_byte)

}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("index path:" + "../web/dist/" + r.URL.Path[1:])
	http.ServeFile(w, r, "../web/dist/"+r.URL.Path[1:])
}

func productsImages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	resp := &model.ImageResponse{}
	err := xclient.Call(context.Background(), "Get", model.ImageRequest(name), resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("productsImages:", model.ImageRequest(name), " len:", resp.ContentLength)

	h := w.Header()
	h.Set("Context-Type", resp.ContentType)
	h.Set("Context-Length", strconv.Itoa(resp.ContentLength))
	w.Write(resp.Content)
	//http.ServeFile(w, r, "../service/static/"+ps.ByName("name"))
}
