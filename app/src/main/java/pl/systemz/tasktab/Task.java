package pl.systemz.tasktab;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.NotificationCompat;
import androidx.core.app.NotificationManagerCompat;

import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.media.RingtoneManager;
import android.net.Uri;
import android.os.Bundle;
import android.provider.Settings;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import static java.lang.System.currentTimeMillis;

public class Task extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_task);

        final Integer taskId = getIntent().getIntExtra("TASK_ID", 0);
//        Log.v("tasktab",taskId.toString());

        EditText editText = (EditText) findViewById(R.id.taskIdInput);
        editText.setText(taskId.toString(), TextView.BufferType.EDITABLE);

        final Button startButton = (Button) findViewById(R.id.taskCounterStart);
        final Button stopButton = (Button) findViewById(R.id.taskCounterStop);
        stopButton.setVisibility(View.GONE);

        startButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startButton.setVisibility(View.GONE);
                stopButton.setVisibility(View.VISIBLE);
                counterStart(taskId);
            }
        });

        stopButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startButton.setVisibility(View.VISIBLE);
                stopButton.setVisibility(View.GONE);
                counterStop(taskId);
            }
        });

    }

    protected void counterStart(Integer taskId) {

//        Bitmap bigIcon = BitmapFactory.decodeResource(getResources(), R.drawable.ic_access_time_black_24dp);
        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, "tasktab")
                .setSmallIcon(R.drawable.ic_access_time_black_24dp)
//                .setLargeIcon(bigIcon)
                .setContentTitle("Task "+taskId.toString())
                .setContentText("Counting time...")
                .setUsesChronometer(true)
                .setWhen(currentTimeMillis())
                .setOngoing(true)
                //.setVibrate(new long[] { 1000, 1000})
                //.setSound(Settings.System.DEFAULT_NOTIFICATION_URI)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT);

//        Uri alarmSound = RingtoneManager.getDefaultUri(RingtoneManager.TYPE_NOTIFICATION);
//        builder.setSound(alarmSound);

        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);

        // notificationId is a unique int for each notification that you must define
        notificationManager.notify(taskId, builder.build());
    }

    protected void counterStop(Integer taskId) {
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
        notificationManager.cancel(taskId);
    }
}
