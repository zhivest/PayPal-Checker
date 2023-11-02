package gg.dragonfruit.paypal;

import java.io.IOException;
import java.net.MalformedURLException;

import org.json.JSONObject;

import gg.dragonfruit.paypal.http.HttpsClient;
import gg.dragonfruit.paypal.http.MimeType;

public class PayPalConnector {

    public Payment getPayment(String paymentId)
            throws MalformedURLException, IOException {
        String response = doGetCall(getUrl("/v1/payments/payment/" + paymentId));
        return response == null ? null : new Payment(new JSONObject(response));
    }

    private String doGetCall(String url)
            throws IOException, MalformedURLException {
        return new HttpsClient(url, HttpsClient.Mode.GET).setAcceptLanguage("en_US")
                .setAccept(MimeType.JSON)
                .setUserPasswordAuthorization(getClientId() + ":" + getClientSecret()).send();
    }

    private String getUrl(String url) {
        return getMode() == Mode.SANDBOX ? "https://api.sandbox.paypal.com" + url
                : "https://api.paypal.com" + url;
    }

    public enum Mode {
        LIVE, SANDBOX
    }

    String clientId, clientSecret;
    Mode mode = Mode.SANDBOX;

    public PayPalConnector(String clientId, String clientSecret, Mode mode) {
        this.clientId = clientId;
        this.clientSecret = clientSecret;
        this.mode = mode;
    }

    public String getClientId() {
        return clientId;
    }

    public String getClientSecret() {
        return clientSecret;
    }

    public Mode getMode() {
        return mode;
    }
}
