package gg.dragonfruit.paypal.http;

import javax.net.ssl.HttpsURLConnection;
import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Base64;

/**
 * Handles all the https requests including the pooling and re-usage of
 * connections.
 */
public class HttpsClient {
    public enum Mode {
        POST, GET
    }

    private String url;
    private Mode mode;
    private String authorization;
    private String acceptLanguage;
    private MimeType contentType;
    private MimeType accept;
    private boolean keepAlive;

    /**
     * It doesn't matter if you re-use this class for same urls or if you create new
     * objects of this class.
     * The underlaying {@link URL#openConnection()} does the connection pooling and
     * re-usage for us.
     * 
     * @param url  Url for the request, e. g.
     *             "https://api.sandbox.paypal.com/v1/payments/payment".
     * @param mode {@link Mode#POST} {@link Mode#GET}
     */
    public HttpsClient(String url, Mode mode) {
        this.url = url;
        this.mode = mode;

    }

    /**
     * Executes the https call. Use this method only for GET calls (GET calls don't
     * have a body).
     * 
     * @return The server's response.
     * @throws IOException If any IOException occurs while connecting the server.
     */
    public String send() throws IOException {
        return send(null);
    }

    /**
     * Executes the https call. Use this method only for POST calls (POST calls need
     * a body).
     * 
     * @param body The body as post input.
     * @return The server's response.
     * @throws IOException If any IOException occurs while connecting the server.
     */
    public String send(String body) throws IOException {
        if (body == null && mode == Mode.POST) {
            throw new IllegalArgumentException("body can't be null for POST calls.");
        } else if (body != null && mode == Mode.GET) {
            throw new IllegalArgumentException("body must be null for GET calls.");
        }
        HttpURLConnection conn = getConnection();
        // if (log.isDebugEnabled()) log.debug(mode.name() + " call '" + url + "' with
        // body: " + body);
        if (mode == Mode.POST) {
            byte[] data = body.getBytes("UTF-8");
            conn.setRequestProperty("Content-Length", String.valueOf(data.length));
            conn.setDoOutput(true);
            try (OutputStream os = conn.getOutputStream()) {
                os.write(data);
            }
        }
        if (conn.getResponseCode() != HttpURLConnection.HTTP_CREATED &&
                conn.getResponseCode() != HttpsURLConnection.HTTP_OK) {
            throw new RuntimeException("Failed : HTTP error code : "
                    + conn.getResponseCode());
        }
        StringWriter out = new StringWriter();
        try (InputStreamReader is = new InputStreamReader(conn.getInputStream())) {
            copy(is, out);
            // not disconnect, leave it to HttpURLConnection
            // conn.disconnect();
            // if (log.isDebugEnabled()) log.debug("Response: " + out.toString());
            return out.toString();
        }
    }

    /**
     * @param accessToken Sets authorization to "Bearer &lt;accessToken&gt;"
     * @return this for chaining.
     */
    public HttpsClient setBearerAuthorization(String accessToken) {
        this.authorization = "Bearer " + accessToken;
        return this;
    }

    /**
     * @param usernamePassword &lt;username&gt;:&lt;password&gt;
     * @return this for chaining.
     */
    public HttpsClient setUserPasswordAuthorization(String usernamePassword) {
        this.authorization = "Basic " + new String(Base64.getEncoder().encode(usernamePassword.getBytes()));
        return this;
    }

    public HttpsClient setAcceptLanguage(String acceptLanguage) {
        this.acceptLanguage = acceptLanguage;
        return this;
    }

    public HttpsClient setContentType(MimeType contentType) {
        this.contentType = contentType;
        return this;
    }

    public HttpsClient setAccept(MimeType accept) {
        this.accept = accept;
        return this;
    }

    public HttpsClient setKeepAlive(boolean keepAlive) {
        this.keepAlive = keepAlive;
        return this;
    }

    private HttpURLConnection getConnection() throws IOException {
        URL turl = new URL(url);
        HttpsURLConnection conn = (HttpsURLConnection) turl.openConnection();
        conn.setRequestMethod(mode == Mode.GET ? "GET" : "POST");
        conn.setDefaultUseCaches(false);
        if (this.keepAlive) {
            conn.setRequestProperty("keep-alive", "true");
        }
        if (authorization != null) {
            conn.setRequestProperty("Authorization", authorization);
        }
        if (acceptLanguage != null) {
            conn.setRequestProperty("Accept-Language", acceptLanguage);
        }
        if (contentType != null) {
            conn.setRequestProperty("Content-Type", contentType.asString());
        }
        if (accept != null) {
            conn.setRequestProperty("Accept", accept.asString());
        }
        return conn;
    }

    private static final int BUFFER_SIZE = 4 * 1024;

    private void copy(final Reader input, final Writer output) throws IOException {
        char[] buffer = new char[BUFFER_SIZE];
        int n;
        while (-1 != (n = input.read(buffer))) {
            output.write(buffer, 0, n);
        }
    }
}
