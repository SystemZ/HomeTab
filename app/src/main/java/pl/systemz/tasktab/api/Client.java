package pl.systemz.tasktab.api;

import java.util.List;

import retrofit2.Call;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;
import retrofit2.http.GET;
import retrofit2.http.POST;
import retrofit2.http.Path;

public class Client {
    public static final String API_URL = "http://192.168.2.88:3000/api/v1/";
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

        public Timer(int id, String name, List<String> tags, int seconds) {
            this.id = id;
            this.name = name;
            this.tags = tags;
            this.seconds = seconds;
        }
    }

    public interface GitHub {
        @GET("/repos/{owner}/{repo}/contributors")
        Call<List<Contributor>> contributors(
                @Path("owner") String owner,
                @Path("repo") String repo);

        @GET("timer")
        Call<List<Timer>> timers();

        @GET("timer/{id}")
        Call<Timer> timerInfo(
                @Path("id") int id
        );

        @POST("timer/{id}/start")
        Call<Timer> timerStart(
                @Path("id") int id
        );

        @POST("timer/{id}/stop")
        Call<Timer> timerStop(
                @Path("id") int id
        );
    }

    private Client() {
        buildRetrofit(API_URL);
    }

    public static Client getInstance() {
        if (instance == null) {
            instance = new Client();
        }
        return instance;
    }

    private void buildRetrofit(String url) {
        Retrofit retrofit = new Retrofit.Builder()
                .baseUrl(url)
                .addConverterFactory(GsonConverterFactory.create())
                .build();

        this.github = retrofit.create(GitHub.class);
    }

    public GitHub getGithub() {
        return this.github;
    }
}