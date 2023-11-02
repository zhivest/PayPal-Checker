package gg.dragonfruit.paypal;

public enum State {
    DENIED("D"), PENDING("P"), SUCCESS("S"), REFUNDED("V");

    final String string;

    State(String string) {
        this.string = string;
    }

    public String getString() {
        return string;
    }

    public static State fromString(String string) {
        for (State state : State.values()) {
            if (state.getString().equalsIgnoreCase(string)) {
                return state;
            }
        }

        return null;
    }
}
