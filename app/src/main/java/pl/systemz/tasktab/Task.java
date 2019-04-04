package pl.systemz.tasktab;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.NavUtils;
import androidx.core.app.NotificationCompat;
import androidx.core.app.NotificationManagerCompat;
import pl.systemz.tasktab.api.Client;
import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

import android.app.PendingIntent;
import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.media.RingtoneManager;
import android.net.Uri;
import android.os.Bundle;
import android.provider.Settings;
import android.util.Log;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import java.io.IOException;
import java.util.List;

import static java.lang.System.currentTimeMillis;

public class Task extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_task);

//        final Integer taskId = getIntent().getIntExtra("TASK_ID", 0);
        final Integer taskIdz = getIntent().getIntExtra("TASK_ID", 0);
        Log.v("task_id", taskIdz.toString());
//        Log.v("tasktab",taskId.toString());

        EditText editText = (EditText) findViewById(R.id.taskIdInput);
        editText.setText(taskIdz.toString(), TextView.BufferType.EDITABLE);

        final Button startButton = (Button) findViewById(R.id.taskCounterStart);
        final Button stopButton = (Button) findViewById(R.id.taskCounterStop);
        stopButton.setVisibility(View.GONE);

        startButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startButton.setVisibility(View.GONE);
                stopButton.setVisibility(View.VISIBLE);
                counterStart(taskIdz);
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

    protected void counterStart(Integer taskId) {

        // this allows returning to task after tapping on notification
        Intent activityIntent = new Intent(this, Task.class);
        activityIntent.putExtra("TASK_ID", taskId);
        Log.v("task_id", taskId.toString());
//        PendingIntent contentIntent = PendingIntent.getActivity(this, 0, activityIntent, PendingIntent.FLAG_UPDATE_CURRENT);
        PendingIntent contentIntent = PendingIntent.getActivity(this, 0, activityIntent, 0);

        Intent broadcastIntent = new Intent(this, NotificationReceiver.class);
        broadcastIntent.putExtra("TASK_ID", taskId);
        PendingIntent actionIntent = PendingIntent.getBroadcast(this, 0, broadcastIntent, PendingIntent.FLAG_UPDATE_CURRENT);
//        PendingIntent.getBroadcast()

//        Bitmap bigIcon = BitmapFactory.decodeResource(getResources(), R.drawable.ic_access_time_black_24dp);
        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, "tasktab-task-" + taskId.toString())
                .setSmallIcon(R.drawable.ic_access_time_black_24dp)
//                .setLargeIcon(bigIcon)
                .setContentTitle("Task " + taskId.toString())
                .setContentText("Counting time...")
                .setUsesChronometer(true)
                .setWhen(currentTimeMillis())
                .setOngoing(true)
                .setContentIntent(actionIntent)
                .addAction(R.mipmap.ic_launcher, "Test", actionIntent)
                .setColor(Color.GREEN)
                //.setAutoCancel(true)
                //.setVibrate(new long[] { 1000, 1000})
                //.setSound(Settings.System.DEFAULT_NOTIFICATION_URI)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT);


//        Uri alarmSound = RingtoneManager.getDefaultUri(RingtoneManager.TYPE_NOTIFICATION);
//        builder.setSound(alarmSound);

        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);

        // notificationId is a unique int for each notification that you must define
        notificationManager.notify(taskId, builder.build());

//        Client client = Client.getInstance();
//        Call<List<Client.Contributor>> call = client.getGithub().contributors("square", "retrofit");
//        call.enqueue(new Callback<List<Client.Contributor>>() {
//            @Override
//            public void onResponse(Call<List<Client.Contributor>> call, Response<List<Client.Contributor>> response) {
//                //System.out.println(response.body().toString());
//                for(Client.Contributor contributor : response.body()) {
//                    System.out.println(contributor.contributions);
//                }
//            }
//
//            @Override
//            public void onFailure(Call<List<Client.Contributor>> call, Throwable t) {
//
//            }
//        });

    }

    protected void counterStop(Integer taskId) {
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
        notificationManager.cancel(taskId);
    }
}
