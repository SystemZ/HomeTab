package pl.systemz.tasktab.api;

import android.content.Context;
import android.content.SharedPreferences;
import android.preference.PreferenceManager;

import java.io.IOException;
import java.util.List;

import okhttp3.Interceptor;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import retrofit2.Call;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;
import retrofit2.http.Body;
import retrofit2.http.GET;
import retrofit2.http.POST;
import retrofit2.http.Path;

public class Client {
    Context context;
    public static final String API_URL = "https://tasktab.lvlup.pro/api/v1/";
    //public static final String API_URL = "http://192.168.2.88:3000/api/v1/";

    private static Client instance = null;
    private GitHub github;

    public static class Contributor {
        public final String login;
        public final int contributions;

        public Contributor(String login, int contributions) {
            this.login = login;
            this.contributions = contributions;
        }
    }

    public static class Timer {
        public final int id;
        public final String name;
        public final List<String> tags;
        public final int seconds;
        public final boolean inProgress;

        public Timer(int id, String name, List<String> tags, int seconds, boolean inProgress) {
            this.id = id;
            this.name = name;
            this.tags = tags;
            this.seconds = seconds;
            this.inProgress = inProgress;
        }
    }

    public static class MqCredentials {
        public final String id;
        public final String host;
        public final int port;
        public final String username;
        public final String password;

        public MqCredentials(String id, String host, int port, String username, String password) {
            this.id = id;
            this.host = host;
            this.port = port;
            this.username = username;
            this.password = password;
        }
    }

    public static class DeviceToken {
        public final String token;

        public DeviceToken(String token) {
            this.token = token;
        }
    }

    public static class PushRegisterRequest {
        final String pushToken;

        public PushRegisterRequest(String pushToken) {
            this.pushToken = pushToken;
        }
    }

    public interface GitHub {
        @GET("/repos/{owner}/{repo}/contributors")
        Call<List<Contributor>> contributors(
                @Path("owner") String owner,
                @Path("repo") String repo);

        @GET("counter")
        Call<List<Timer>> timers();

        @GET("counter/{id}")
        Call<Timer> timerInfo(
                @Path("id") int id
        );

        @POST("counter/{id}/start")
        Call<Timer> timerStart(
                @Path("id") int id
        );

        @POST("counter/{id}/stop")
        Call<Timer> timerStop(
                @Path("id") int id
        );

        // TODO delete this
        @GET("mq/access")
        Call<MqCredentials> mqCredentialsGet();

        @POST("push/register")
        Call<Void> deviceRegister(
                @Body PushRegisterRequest body
        );

    }

    private Client(Context context) {
        this.context = context;
        buildRetrofit(API_URL);
    }

    public static Client getInstance(Context context) {
        if (instance == null) {
            instance = new Client(context);
        }
        return instance;
    }

    private void buildRetrofit(String url) {

        SharedPreferences prefs = PreferenceManager.getDefaultSharedPreferences(context);
        final String token = prefs.getString("auth_token", "emptytoken");

        OkHttpClient client = new OkHttpClient.Builder().addInterceptor(new Interceptor() {
            @Override
            public Response intercept(Chain chain) throws IOException {
                Request newRequest = chain.request().newBuilder()
                        .addHeader("Authorization", token)
//                        .addHeader("Authorization", "Bearer " + token)
                        .build();
                return chain.proceed(newRequest);
            }
        }).build();

        Retrofit retrofit = new Retrofit.Builder()
                .client(client)
                .baseUrl(url)
                .addConverterFactory(GsonConverterFactory.create())
                .build();

        this.github = retrofit.create(GitHub.class);
    }

    public GitHub getGithub() {
        return this.github;
    }
}