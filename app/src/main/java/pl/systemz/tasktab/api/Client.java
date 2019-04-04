package pl.systemz.tasktab.api;

import java.util.List;

import retrofit2.Call;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;
import retrofit2.http.GET;
import retrofit2.http.Path;

public class Client {
    public static final String API_URL = "https://api.github.com";
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

    public interface GitHub {
        @GET("/repos/{owner}/{repo}/contributors")
        Call<List<Contributor>> contributors(
                @Path("owner") String owner,
                @Path("repo") String repo);
    }

    private Client() {
        buildRetrofit(API_URL);
    }

    public static Client getInstance() {
        if(instance == null) {
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