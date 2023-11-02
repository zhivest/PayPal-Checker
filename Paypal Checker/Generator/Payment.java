package gg.dragonfruit.paypal;

import org.json.JSONObject;

public class Payment {
    State state;
    float price;

    public Payment(JSONObject paymentObject) {
        JSONObject transactionInfo = paymentObject.getJSONObject("transaction_info");
        this.state = State.fromString(transactionInfo.getString("transaction_status"));
        JSONObject transactionAmount = transactionInfo.getJSONObject("transaction_amount");
        this.price = transactionAmount.getFloat("value") / 100;
    }

    public State getState() {
        return this.state;
    }

    public float getPrice() {
        return this.price;

    }
}
