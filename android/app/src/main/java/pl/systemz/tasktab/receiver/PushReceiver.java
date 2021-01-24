package pl.systemz.tasktab.receiver;

import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.graphics.Color;
import android.os.Build;
import android.provider.Settings;
import android.util.Log;

import androidx.core.app.NotificationCompat;
import androidx.core.app.NotificationManagerCompat;
import pl.systemz.tasktab.R;
import pl.systemz.tasktab.Task;

import static java.lang.System.currentTimeMillis;

public class PushReceiver extends BroadcastReceiver {
    private static final String TAG = "TaskTab";

    @Override
    public void onReceive(Context context, Intent intent) {
        // Prepare a notification with vibration, sound and lights
//        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, "tasktab-service-running")
//                .setSmallIcon(R.drawable.ic_assignment_black_24dp)
//                .setContentTitle("")
//                //.setContentText("content text");
//                .setOngoing(true)
//                .setAutoCancel(false)
//                .setContentIntent(pi);

        /*
        NotificationCompat.Builder builder = new NotificationCompat.Builder(context)
                .setAutoCancel(true)
                .setSmallIcon(android.R.drawable.ic_dialog_info)
                .setContentTitle(notificationTitle)
                .setContentText(notificationText)
                .setLights(Color.RED, 1000, 1000)
                .setVibrate(new long[]{0, 400, 250, 400})
                .setSound(RingtoneManager.getDefaultUri(RingtoneManager.TYPE_NOTIFICATION))
                .setContentIntent(PendingIntent.getActivity(context, 0, new Intent(context, MainActivity.class), PendingIntent.FLAG_UPDATE_CURRENT));

        // Automatically configure a Notification Channel for devices running Android O+
        Pushy.setNotificationChannel(builder, context);
        // Get an instance of the NotificationManager service
        NotificationManager notificationManager = (NotificationManager) context.getSystemService(context.NOTIFICATION_SERVICE);
        // Build the notification and display it
        notificationManager.notify(1, builder.build());
        */

        handlePush(getResultData(), context, intent);
    }

    public void handlePush(String payload, Context context, Intent intent) {
        Log.d(TAG, "Got new push msg");
        // action based on msg
        if (intent.getStringExtra("type").equals("startNotification")) {
            Log.d(TAG, "Creating new notification...");
            startOngoingNotification(intent, context);
        } else if (intent.getStringExtra("type").equals("stopNotification")) {
            Log.d(TAG, "Removing notification...");
            stopOngoingNotification(intent, context);
        } else if (intent.getStringExtra("type").equals("showNotification")) {
            Log.d(TAG, "Showing ad-hoc notification...");
            showNotification(intent, context);
        } else {
            Log.d(TAG, "Task in message not recognized...");
        }
    }

    private void stopOngoingNotification(Intent intent, Context context) {
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(context);
        notificationManager.cancel(intent.getIntExtra("sessionId", 0));
    }

    private void startOngoingNotification(Intent intent, Context context) {
        // prepare handler for clicking notification
        Intent taskView = new Intent(context, Task.class);
        taskView.putExtra("TASK_ID", intent.getIntExtra("id", 0));
        taskView.putExtra("TASK_NAME", "...");
        // use System.currentTimeMillis() to have a unique ID for the pending intent
        PendingIntent pIntent = PendingIntent.getActivity(context, (int) System.currentTimeMillis(), taskView, 0);

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel = new NotificationChannel("tasktab-counters",
                    "Active counters",
                    NotificationManager.IMPORTANCE_DEFAULT);
            NotificationManagerCompat notificationManager = NotificationManagerCompat.from(context);
            notificationManager.createNotificationChannel(channel);
        }
        NotificationCompat.Builder builder = new NotificationCompat.Builder(context, "tasktab-counters")
                .setSmallIcon(R.drawable.ic_access_time_black_24dp)
                .setContentTitle(intent.getStringExtra("title"))
                .setContentText(intent.getStringExtra("msg"))
                .setUsesChronometer(true)
                .setVisibility(NotificationCompat.VISIBILITY_PUBLIC)
                .setWhen(currentTimeMillis())
                .setOngoing(true)
                .setContentIntent(pIntent)
                //.addAction(R.mipmap.ic_launcher, "Stop", stopActionIntent)
                .setColor(Color.GREEN)
                .setSound(Settings.System.DEFAULT_NOTIFICATION_URI)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT);
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(context);
        notificationManager.notify(intent.getIntExtra("sessionId", 0), builder.build());
    }

    private void showNotification(Intent intent, Context context) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel = new NotificationChannel("tasktab-generic-notifications",
                    "Generic notify",
                    NotificationManager.IMPORTANCE_DEFAULT);
            NotificationManagerCompat notificationManager = NotificationManagerCompat.from(context);
            notificationManager.createNotificationChannel(channel);
        }
        NotificationCompat.Builder builder = new NotificationCompat.Builder(context, "tasktab-counters")
                .setSmallIcon(R.drawable.ic_assignment_black_24dp)
                .setContentTitle(intent.getStringExtra("title"))
                .setContentText(intent.getStringExtra("msg"))
                .setVisibility(NotificationCompat.VISIBILITY_PUBLIC)
                .setSound(Settings.System.DEFAULT_NOTIFICATION_URI)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT);
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(context);
        // use System.currentTimeMillis() to have a unique ID for notification
        notificationManager.notify((int) System.currentTimeMillis(), builder.build());
    }

}