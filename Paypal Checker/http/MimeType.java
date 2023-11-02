package gg.dragonfruit.paypal.http;

public enum MimeType {
    JSON("application/json");

    private String str;

    MimeType(String str) {
        this.str = str;
    }

    public String asString() {
        return str;
    }
}
