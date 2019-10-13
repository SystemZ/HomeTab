package pl.systemz.tasktab;

import android.app.AlarmManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Intent;
import android.content.IntentFilter;
import android.graphics.Color;
import android.os.IBinder;
import android.util.Log;
import android.widget.Toast;

import com.google.gson.Gson;
import com.hivemq.client.mqtt.MqttClient;
import com.hivemq.client.mqtt.mqtt3.Mqtt3AsyncClient;
import com.hivemq.client.mqtt.mqtt3.message.connect.connack.Mqtt3ConnAck;
import com.hivemq.client.mqtt.mqtt3.message.publish.Mqtt3Publish;
import com.hivemq.client.mqtt.mqtt3.message.subscribe.suback.Mqtt3SubAck;

import java.nio.charset.StandardCharsets;

import androidx.core.app.NotificationCompat;
import androidx.core.app.NotificationManagerCompat;
import pl.systemz.tasktab.model.MqttMsg;

import static java.lang.System.currentTimeMillis;

public class StalkService extends Service {
    private static final String TAG = "Stalk";
    private Antenna antenna;

    private AlarmManager alarmMgr;
    private PendingIntent alarmIntent;

    int mStartMode;       // indicates how to behave if the service is killed
    IBinder mBinder;      // interface for clients that bind
    boolean mAllowRebind; // indicates whether onRebind should be used

    Mqtt3AsyncClient mqClient;

    @Override
    public void onCreate() {
        // The service is being created

        //Intent intents = new Intent(getBaseContext(),MainActivity.class);
        //intents.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        //startActivity(intents);

        Toast.makeText(this, "Stalk service started", Toast.LENGTH_LONG).show();
        Log.d(TAG, "onStart, register events");

        // Hopefully your alarm will have a lower frequency than this!
        // AlarmManager alarmMgr;
        // alarmMgr = (AlarmManager)this.getSystemService(Context.ALARM_SERVICE);
        // pi = PendingIntent.getBroadcast(this, 11, new Intent("pl.systemz.stalk.GPSSETTINGS"), 0);
        // am.setInexactRepeating(AlarmManager.RTC_WAKEUP, 0, TWENTYFIVE_SECONDS, pi);

        // alarmMgr.setInexactRepeating(AlarmManager.ELAPSED_REALTIME_WAKEUP,
        //         SystemClock.elapsedRealtime() + 60*1000,
        //         60*1000, alarmIntent);

        // SystemClock.elapsedRealtime() + AlarmManager.INTERVAL_HALF_HOUR,
        // AlarmManager.INTERVAL_HALF_HOUR, alarmIntent);

        IntentFilter intentFilter = new IntentFilter();
        intentFilter.addAction(Intent.ACTION_SCREEN_ON);
        intentFilter.addAction(Intent.ACTION_SCREEN_OFF);
        intentFilter.addAction(Intent.ACTION_BATTERY_CHANGED);
        intentFilter.addAction(Intent.ACTION_BOOT_COMPLETED);

        this.antenna = new Antenna();
        registerReceiver(antenna, intentFilter);

        // connect to MQTT server
        mqClient = MqttClient
                .builder()
                .useMqttVersion3()
                .identifier("tasktab-android")
                .serverHost("changeme")
                .serverPort(1883)
                .automaticReconnect()
                .applyAutomaticReconnect()
                .addDisconnectedListener(disconnectedContext -> {
                    Log.d(TAG, "Got MQTT DC");
                    disconnectedContext.getReconnector().reconnect(true);
                })
                .addConnectedListener(connectedContext -> {
                    Log.d(TAG, "MQTT connected");
                })
                //.useSslWithDefaultConfig()
                .buildAsync();

        // login to MQTT server
        mqClient.connectWith()
                .simpleAuth()
                .username("tasktab:tasktab-android")
                .password("changeme".getBytes())
                .applySimpleAuth()
                .send()
                .whenComplete(this::mqConnectionComplete);
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        // The service is starting, due to a call to startService()
        Log.d(TAG, "onStartCommand");

        return mStartMode;
    }

    @Override
    public IBinder onBind(Intent intent) {
        // A mqClient is binding to the service with bindService()
        return mBinder;
    }

    @Override
    public boolean onUnbind(Intent intent) {
        // All clients have unbound with unbindService()
        return mAllowRebind;
    }

    @Override
    public void onRebind(Intent intent) {
        // A mqClient is binding to the service with bindService(),
        // after onUnbind() has already been called
    }

    @Override
    public void onDestroy() {
        // The service is no longer used and is being destroyed
        Toast.makeText(this, "Stalk Service stopped", Toast.LENGTH_LONG).show();
        Log.d(TAG, "onDestroy");
    }

    // MQTT related below
    //
    private void mqConnectionComplete(Mqtt3ConnAck connAck, Throwable throwable) {
        if (throwable != null) {
            // handle failure
            Log.d(TAG, "Failure connecting to MQTT server");
        } else {
            // setup subscribes or start publishing
            Log.d(TAG, "Connected to MQTT server");
            mqClient.subscribeWith()
                    .topicFilter("tasktab")
                    // Process the received message
                    .callback(this::newMqMsg)
                    .send()
                    .whenComplete(StalkService::subscribeComplete);
        }
    }

    private static void subscribeComplete(Mqtt3SubAck subAck, Throwable throwable2) {
        if (throwable2 != null) {
            // Handle failure to subscribe
            Log.d(TAG, "Subscribe error");
        } else {
            // Handle successful subscription, e.g. logging or incrementing a metric
            Log.d(TAG, "Subscribe OK");
        }
    }

    private void newMqMsg(Mqtt3Publish publish) {
        Log.d(TAG, "Got a new msg via MQTT");
        String msgTxt = new String(publish.getPayloadAsBytes(), StandardCharsets.UTF_8);
        // deserialize
        // {"type":"startNotification","msg":"testzz",id:1}
        // {"type":"stopNotification","msg":"testzz",id:1}
        Gson gson = new Gson();
        MqttMsg msgObj = gson.fromJson(msgTxt, MqttMsg.class);
        // action based on msg
        if (msgObj.getType().equals("startNotification")) {
            startNotification(msgObj);
        } else if (msgObj.getType().equals("stopNotification")) {
            stopNotification(msgObj);
        }
    }

    private void stopNotification(MqttMsg msgObj) {
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
        notificationManager.cancel(msgObj.getId());
    }

    private void startNotification(MqttMsg msgObj) {
        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, "tasktab-counters")
                .setSmallIcon(R.drawable.ic_access_time_black_24dp)
                .setContentTitle(msgObj.getMsg())
                .setContentText("In progress...")
                .setUsesChronometer(true)
                .setWhen(currentTimeMillis())
                //.setOngoing(true)
                //.setContentIntent(tapIntent)
                //.addAction(R.mipmap.ic_launcher, "Stop", stopActionIntent)
                .setColor(Color.GREEN)
                //.setAutoCancel(true)
                .setVibrate(new long[]{150})
                //.setSound(Settings.System.DEFAULT_NOTIFICATION_URI)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT);
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
        // notificationId is a unique int for each notification that you must define
        notificationManager.notify(msgObj.getId(), builder.build());
    }

}
