package pl.systemz.tasktab;

import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuInflater;
import android.view.MenuItem;

import java.util.ArrayList;
import java.util.List;

import androidx.appcompat.app.AppCompatActivity;
import androidx.appcompat.widget.SearchView;
import androidx.core.app.NavUtils;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;
import pl.systemz.tasktab.api.Client;
import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

//import android.widget.SearchView;

public class TaskList extends AppCompatActivity {
    private RecyclerView recyclerView;
    private TaskListAdapter mAdapter;
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
        final List<Client.Timer> input = new ArrayList<>();

        //FIXME need progress bar
        Client.Timer loadingInProgress = new Client.Timer(0, "Loading...", new ArrayList<String>(), 0, false);
        input.add(loadingInProgress);
        mAdapter = new TaskListAdapter(input);
        recyclerView.setAdapter(mAdapter);

        // calling backend API
        Client client = Client.getInstance();
        Call<List<Client.Timer>> call = client.getGithub().timers();
        call.enqueue(new Callback<List<Client.Timer>>() {
            @Override
            public void onResponse(Call<List<Client.Timer>> call, Response<List<Client.Timer>> response) {
                if (!response.isSuccessful()) {
                    return;
                }
                // remove loading task
                input.remove(0);
                for (Client.Timer timer : response.body()) {
                    Client.Timer task = new Client.Timer(timer.id, timer.name, timer.tags, timer.seconds, timer.inProgress);
                    input.add(task);
                }
                // define an adapter
                mAdapter = new TaskListAdapter(input);
                recyclerView.setAdapter(mAdapter);
            }

            @Override
            public void onFailure(Call<List<Client.Timer>> call, Throwable t) {
                // remove loading task
                input.remove(0);
                //FIXME
                Client.Timer failure = new Client.Timer(0, "Loading tasks failed :(", new ArrayList<String>(), 0, false);
                input.add(failure);
                // define an adapter
                mAdapter = new TaskListAdapter(input);
                recyclerView.setAdapter(mAdapter);
            }
        });

    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        MenuInflater inflater = getMenuInflater();
        inflater.inflate(R.menu.menu_task_list, menu);

        MenuItem searchItem = menu.findItem(R.id.action_search);
        SearchView searchView = (SearchView) searchItem.getActionView();
        searchView.setOnQueryTextListener(new SearchView.OnQueryTextListener() {
            @Override
            public boolean onQueryTextSubmit(String query) {
                Log.v("", "onQueryTextSubmit");
                return false;
            }

            @Override
            public boolean onQueryTextChange(String newText) {
                // https://www.youtube.com/watch?v=sJ-Z9G0SDhc
                mAdapter.getFilter().filter(newText);
                Log.v("", "onQueryTextChange");
                return false;
            }
        });

        return true;
    }
//

    public boolean onQueryTextListener(String query) {
        // Here is where we are going to implement the filter logic
        return false;
    }
//
//    @Override
//    public boolean onQueryTextSubmit(String query) {
//        return false;
//    }

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
