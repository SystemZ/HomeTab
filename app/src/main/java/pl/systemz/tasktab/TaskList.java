package pl.systemz.tasktab;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.NavUtils;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;
import pl.systemz.tasktab.api.Client;
import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

import android.os.Bundle;
import android.util.Log;
import android.view.MenuItem;
import android.widget.Toast;

import java.util.ArrayList;
import java.util.List;

public class TaskList extends AppCompatActivity {
    private RecyclerView recyclerView;
    private RecyclerView.Adapter mAdapter;
    private RecyclerView.LayoutManager layoutManager;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_task_list);

        recyclerView = findViewById(R.id.my_recycler_view);
        // use this setting to
        // improve performance if you know that changes
        // in content do not change the layout size
        // of the RecyclerView
        recyclerView.setHasFixedSize(true);
        // use a linear layout manager
        layoutManager = new LinearLayoutManager(this);
        recyclerView.setLayoutManager(layoutManager);
        final List<String> input = new ArrayList<>();
        input.add("Loading...");

        Client client = Client.getInstance();
        Call<List<Client.Timer>> call = client.getGithub().timers();
        call.enqueue(new Callback<List<Client.Timer>>() {
            @Override
            public void onResponse(Call<List<Client.Timer>> call, Response<List<Client.Timer>> response) {
                if (!response.isSuccessful()) {
                    return;
                }
                input.remove(0);
                for (Client.Timer timer : response.body()) {
//                    System.out.println(timer.id);
//                    System.out.println(timer.name);
                    input.add(timer.name);
                }
                // define an adapter
                mAdapter = new TaskListAdapter(input);
                recyclerView.setAdapter(mAdapter);
            }

            @Override
            public void onFailure(Call<List<Client.Timer>> call, Throwable t) {
                input.remove(0);
                input.add("Loading tasks failed");
                // define an adapter
                mAdapter = new TaskListAdapter(input);
                recyclerView.setAdapter(mAdapter);
            }
        });

    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        Log.v("menu", item.toString());
        switch (item.getItemId()) {
            // Respond to the action bar's Up/Home button
            case android.R.id.home:
                NavUtils.navigateUpFromSameTask(this);
//                Intent mainActivityIntent = new Intent(this, MainActivity.class);
//                this.startActivity(mainActivityIntent);
                return true;
        }
        return super.onOptionsItemSelected(item);
    }
//
//    @Override
//    public void onBackPressed() {
//        this.finish();
//        Intent mainActivityIntent = new Intent(this, MainActivity.class);
//        this.startActivity(mainActivityIntent);
//    }

}
