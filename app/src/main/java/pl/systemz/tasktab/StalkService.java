package pl.systemz.tasktab;

import android.app.AlarmManager;
import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Intent;
import android.graphics.Color;
import android.os.Build;
import android.os.IBinder;
import android.provider.Settings;
import android.util.Log;
import android.widget.Toast;

import com.google.gson.Gson;
import com.hivemq.client.mqtt.MqttClient;
import com.hivemq.client.mqtt.MqttClientState;
import com.hivemq.client.mqtt.mqtt3.Mqtt3AsyncClient;
import com.hivemq.client.mqtt.mqtt3.message.connect.connack.Mqtt3ConnAck;
import com.hivemq.client.mqtt.mqtt3.message.publish.Mqtt3Publish;
import com.hivemq.client.mqtt.mqtt3.message.subscribe.suback.Mqtt3SubAck;

import java.nio.charset.StandardCharsets;
import java.util.concurrent.TimeUnit;

import androidx.core.app.NotificationCompat;
import androidx.core.app.NotificationManagerCompat;
import pl.systemz.tasktab.api.Client;
import pl.systemz.tasktab.model.MqttMsg;
import pl.systemz.tasktab.receiver.Antenna;
import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

import static java.lang.System.currentTimeMillis;

public class StalkService extends Service {
    private static final String TAG = "StalkService";
    private Antenna antenna;

    private AlarmManager alarmMgr;
    private PendingIntent alarmIntent;

    IBinder mBinder;      // interface for clients that bind
    boolean mAllowRebind; // indicates whether onRebind should be used

    Mqtt3AsyncClient mqClient;

    @Override
    public void onCreate() {
        // The service is being created
        Log.d(TAG, "onCreate()");

        //Intent intents = new Intent(getBaseContext(),MainActivity.class);
        //intents.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        //startActivity(intents);

        Toast.makeText(this, "Stalk service started", Toast.LENGTH_LONG).show();

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

//        IntentFilter intentFilter = new IntentFilter();
//        intentFilter.addAction(Intent.ACTION_SCREEN_ON);
//        intentFilter.addAction(Intent.ACTION_SCREEN_OFF);
//        intentFilter.addAction(Intent.ACTION_BATTERY_CHANGED);
//        intentFilter.addAction(Intent.ACTION_BOOT_COMPLETED);
//        this.antenna = new Antenna();
//        registerReceiver(antenna, intentFilter);

        // get credentials to MQTT server by asking tasktab backend RESTful API
        Client client = Client.getInstance(getApplicationContext());
        Call<Client.MqCredentials> call = client.getGithub().mqCredentialsGet();
        call.enqueue(new Callback<Client.MqCredentials>() {
            @Override
            public void onResponse(Call<Client.MqCredentials> call, Response<Client.MqCredentials> response) {
                if (!response.isSuccessful()) {
                    // debug
                    Log.e(TAG, "MQ credentials fetch failed");
                    justStop();
                    return;
                }
                Log.d(TAG, "MQ credentials fetched");
                mqConnect(response);
            }

            @Override
            public void onFailure(Call<Client.MqCredentials> call, Throwable t) {
                Log.e(TAG, "MQ credentials fetch onFailure");
                Log.e(TAG, t.getCause().toString());
                justStop();
            }
        });
    }

    // end work if fetching credentials doesn't work
    public void justStop() {
        this.stopSelf();
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        // The service is starting, due to a call to startService()
        Log.d(TAG, "onStartCommand");

        // Show notification to prevent killing in background
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel = new NotificationChannel("tasktab-service-running",
                    "Placeholder, disable me",
                    NotificationManager.IMPORTANCE_DEFAULT);
            NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
            notificationManager.createNotificationChannel(channel);
        }
        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, "tasktab-service-running");

        builder.setSmallIcon(R.drawable.ic_launcher_foreground);
        builder.setContentTitle(".");
        //builder.setContentText("content text");
        //final Intent notificationIntent = new Intent(this, FakeActivity.class);
        //final PendingIntent pi = PendingIntent.getActivity(this, 0, notificationIntent, 0);
        //builder.setContentIntent(pi);
        final Notification notification = builder.build();
        startForeground(1, notification);

        return Service.START_STICKY;
    }

    @Override
    public IBinder onBind(Intent intent) {
        Log.d(TAG, "onBind");
        // A mqClient is binding to the service with bindService()
        return mBinder;
    }

    @Override
    public boolean onUnbind(Intent intent) {
        Log.d(TAG, "onUnbind");
        // All clients have unbound with unbindService()
        return mAllowRebind;
    }

    @Override
    public void onRebind(Intent intent) {
        Log.d(TAG, "onRebind");
        // A mqClient is binding to the service with bindService(),
        // after onUnbind() has already been called
    }

    @Override
    public void onDestroy() {
        // The service is no longer used and is being destroyed
        Toast.makeText(this, "Stalk Service stopped", Toast.LENGTH_LONG).show();
        Log.d(TAG, "onDestroy");
        try {
            mqClient.disconnect();
        } catch (NullPointerException e) {
            Log.d(TAG, "tried to destroy empty mqClient");
        }
    }

    // MQTT related below
    //
    private void mqConnect(Response<Client.MqCredentials> response) {
        Log.d(TAG, "Connecting to MQTT server");
        // connect to MQTT server
        mqClient = MqttClient
                .builder()
                .useMqttVersion3()
                .identifier(response.body().id)
                .serverHost(response.body().host)
                .serverPort(response.body().port)
                .automaticReconnect()
                .initialDelay(1, TimeUnit.SECONDS)
                .maxDelay(10, TimeUnit.SECONDS)
                .applyAutomaticReconnect()
                .addConnectedListener(connectedContext -> {
                    Log.d(TAG, "MQTT connected");
                })
                .addDisconnectedListener(disconnectedContext -> {
                    Log.d(TAG, "Got MQTT DC: " + disconnectedContext.getClientConfig().getState());
                    if (disconnectedContext.getClientConfig().getState() == MqttClientState.CONNECTING_RECONNECT) {
                        disconnectedContext.getReconnector().reconnect(false);
                        this.stopSelf();
                    }
                })
                .buildAsync();

        // login to MQTT server
        mqClient.connectWith()
                .cleanSession(true)
                .keepAlive(10)
                .simpleAuth()
                .username(response.body().username)
                .password(response.body().password.getBytes())
                .applySimpleAuth()
                .send()
                .whenComplete(this::mqConnectionComplete);
        Log.d(TAG, "mqConnect() end");
    }

    private void mqConnectionComplete(Mqtt3ConnAck connAck, Throwable throwable) {
        Log.d(TAG, connAck.toString());
        if (throwable != null) {
            // handle failure
            Log.e(TAG, "Failure connecting to MQTT server");
        } else {
            // setup subscribes or start publishing
            Log.d(TAG, "Connected to MQTT server");
            mqClient.subscribeWith()
                    .topicFilter(mqClient.getConfig().getClientIdentifier().toString())
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
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel = new NotificationChannel("tasktab-counters",
                    "Active counters",
                    NotificationManager.IMPORTANCE_DEFAULT);
            NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
            notificationManager.createNotificationChannel(channel);
        }
        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, "tasktab-counters")
                .setSmallIcon(R.drawable.ic_access_time_black_24dp)
                .setContentTitle(msgObj.getMsg())
                .setContentText("In progress...")
                .setUsesChronometer(true)
                .setWhen(currentTimeMillis())
                .setOngoing(true)
                //.setContentIntent(tapIntent)
                //.addAction(R.mipmap.ic_launcher, "Stop", stopActionIntent)
                .setColor(Color.GREEN)
                .setVibrate(new long[]{150})
                .setSound(Settings.System.DEFAULT_NOTIFICATION_URI)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT);
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
        notificationManager.notify(msgObj.getId(), builder.build());
    }
}
