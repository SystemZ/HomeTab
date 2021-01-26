package pl.systemz.hometab;

import android.app.PendingIntent;
import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;
import pl.systemz.hometab.api.Client;
import pl.systemz.hometab.receiver.NotificationReceiver;
import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class Task extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_task);

        final Integer taskIdz = getIntent().getIntExtra("TASK_ID", 0);
        final String taskName = getIntent().getStringExtra("TASK_NAME");

        final TextView taskTitle = findViewById(R.id.TaskTitle);
        taskTitle.setText(taskName);

        final Button startButton = findViewById(R.id.taskCounterStart);
        final Button stopButton = findViewById(R.id.taskCounterStop);

        Client client = Client.getInstance(getApplicationContext());
        Call<Client.Timer> call = client.getTtClient().timerInfo(taskIdz);
        call.enqueue(new Callback<Client.Timer>() {
            @Override
            public void onResponse(Call<Client.Timer> call, Response<Client.Timer> response) {
                if (!response.isSuccessful()) {
                    Toast.makeText(Task.this, "Something is wrong with the server...", Toast.LENGTH_SHORT).show();
                    return;
                }

                taskTitle.setText(response.body().name);

                if (response.body().inProgress) {
                    stopButton.setVisibility(View.VISIBLE);
                    startButton.setVisibility(View.INVISIBLE);
                } else {
                    stopButton.setVisibility(View.INVISIBLE);
                    startButton.setVisibility(View.VISIBLE);
                }
            }

            @Override
            public void onFailure(Call<Client.Timer> call, Throwable t) {
                Toast.makeText(Task.this, t.getMessage(), Toast.LENGTH_SHORT).show();
            }
        });


        startButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startButton.setVisibility(View.GONE);
                stopButton.setVisibility(View.VISIBLE);
                counterStart(taskIdz, taskName);
            }
        });

        stopButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startButton.setVisibility(View.VISIBLE);
                stopButton.setVisibility(View.GONE);
                counterStop(taskIdz);
            }
        });

    }

    protected void counterStart(Integer taskId, String taskName) {

        // this allows returning to task after tapping on notification
        Intent tapBroadcastIntent = new Intent(this, NotificationReceiver.class);
        tapBroadcastIntent.putExtra("TASK_ID", taskId);
        tapBroadcastIntent.putExtra("TASK_NAME", taskName);
        tapBroadcastIntent.putExtra("ACTION", "VIEW");
        PendingIntent tapIntent = PendingIntent.getBroadcast(this, 0, tapBroadcastIntent, PendingIntent.FLAG_UPDATE_CURRENT);

        // this makes stop button in notification
//        Intent stopActionBroadcastIntent = new Intent(this, NotificationReceiver.class);
//        stopActionBroadcastIntent.putExtra("TASK_ID", taskId);
//        stopActionBroadcastIntent.putExtra("ACTION", "STOP");
//        PendingIntent stopActionIntent = PendingIntent.getBroadcast(this, 1, stopActionBroadcastIntent, PendingIntent.FLAG_UPDATE_CURRENT);


//        Uri alarmSound = RingtoneManager.getDefaultUri(RingtoneManager.TYPE_NOTIFICATION);
//        builder.setSound(alarmSound);

        Client client = Client.getInstance(getApplicationContext());
        Call<Client.Timer> call = client.getTtClient().timerStart(taskId);
        call.enqueue(new Callback<Client.Timer>() {
            @Override
            public void onResponse(Call<Client.Timer> call, Response<Client.Timer> response) {
                if (!response.isSuccessful()) {
                    Toast.makeText(Task.this, "Something is wrong with the server...", Toast.LENGTH_SHORT).show();
                    return;
                }
                Toast.makeText(Task.this, "Counter started", Toast.LENGTH_SHORT).show();
            }

            @Override
            public void onFailure(Call<Client.Timer> call, Throwable t) {
                Toast.makeText(Task.this, t.getMessage(), Toast.LENGTH_SHORT).show();
            }
        });

    }

    protected void counterStop(final Integer taskId) {
        Client client = Client.getInstance(getApplicationContext());
        Call<Client.Timer> call = client.getTtClient().timerStop(taskId);
        call.enqueue(new Callback<Client.Timer>() {
            @Override
            public void onResponse(Call<Client.Timer> call, Response<Client.Timer> response) {
                if (!response.isSuccessful()) {
                    Toast.makeText(Task.this, "Something is wrong with the server...", Toast.LENGTH_SHORT).show();
                    return;
                }
                Toast.makeText(Task.this, "Counter stopped", Toast.LENGTH_SHORT).show();
            }

            @Override
            public void onFailure(Call<Client.Timer> call, Throwable t) {
                Toast.makeText(Task.this, t.getMessage(), Toast.LENGTH_SHORT).show();
            }
        });

    }
}
