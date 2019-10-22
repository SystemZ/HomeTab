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

import org.fusesource.hawtbuf.Buffer;
import org.fusesource.hawtbuf.UTF8Buffer;
import org.fusesource.mqtt.client.CallbackConnection;
import org.fusesource.mqtt.client.Listener;
import org.fusesource.mqtt.client.MQTT;
import org.fusesource.mqtt.client.QoS;
import org.fusesource.mqtt.client.Topic;

import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;

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

    MQTT mqttClient;
    CallbackConnection mqttConnection;

    @Override
    public void onCreate() {
        // The service is being created
        Log.d(TAG, "onCreate()");
        Toast.makeText(this, "Stalk service started", Toast.LENGTH_LONG).show();

        //Intent intents = new Intent(getBaseContext(),MainActivity.class);
        //intents.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        //startActivity(intents);

        // SystemClock.elapsedRealtime() + AlarmManager.INTERVAL_HALF_HOUR,
        // AlarmManager.INTERVAL_HALF_HOUR, alarmIntent);

//        IntentFilter intentFilter = new IntentFilter();
//        intentFilter.addAction(Intent.ACTION_SCREEN_ON);
//        intentFilter.addAction(Intent.ACTION_SCREEN_OFF);
//        intentFilter.addAction(Intent.ACTION_BATTERY_CHANGED);
//        intentFilter.addAction(Intent.ACTION_BOOT_COMPLETED);
//        this.antenna = new Antenna();
//        registerReceiver(antenna, intentFilter);

    }

    // end work if fetching credentials doesn't work
    public void justStop() {
        this.stopSelf();
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        // The service is starting, due to a call to startService()
        Log.d(TAG, "onStartCommand");
        // show notification preventing killing service
        foregroundNotification();

        // get credentials to MQTT server by asking tasktab backend RESTful API
        Client client = Client.getInstance(getApplicationContext());
        Call<Client.MqCredentials> call = client.getGithub().mqCredentialsGet();
        call.enqueue(new Callback<Client.MqCredentials>() {
            @Override
            public void onResponse(Call<Client.MqCredentials> call, Response<Client.MqCredentials> response) {
                if (!response.isSuccessful()) {
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
        mqttConnection.disconnect(new org.fusesource.mqtt.client.Callback<Void>() {
            @Override
            public void onSuccess(Void value) {
                Log.d(TAG, "disconnect from MQTT OK");
            }

            @Override
            public void onFailure(Throwable value) {
                Log.d(TAG, "disconnect from MQTT error");
            }
        });
    }

    private void mqConnect(Response<Client.MqCredentials> response) {
        Log.d(TAG, "Connecting to MQTT server");

        mqttClient = new MQTT();
        try {
            mqttClient.setHost(response.body().host, response.body().port);
        } catch (URISyntaxException e) {
            e.printStackTrace();
        }
        mqttClient.setCleanSession(false);
        mqttClient.setClientId(response.body().id);
        mqttClient.setUserName(response.body().username);
        mqttClient.setPassword(response.body().password);
        mqttClient.setReconnectDelayMax(5000);
        mqttClient.setReconnectBackOffMultiplier(1);
        mqttConnection = mqttClient.callbackConnection();

        mqttConnection.listener(new Listener() {
            public void onDisconnected() {
                Log.d(TAG, "listener onDisconnected");
            }

            public void onConnected() {
                Log.d(TAG, "listener onConnected");
            }

            public void onPublish(UTF8Buffer topic, Buffer payload, Runnable ack) {
                Log.d(TAG, "listener onPublish");
                // You can now process a received message from a topic.
                mqttMsgHandler(payload.toByteArray());
                // Once process execute the ack runnable.
                ack.run();
            }

            public void onFailure(Throwable value) {
                Log.d(TAG, "listener onFailure");
                //mqttConnection.close(null); // a mqttConnection failure occured.
            }
        });

        mqttConnection.connect(new org.fusesource.mqtt.client.Callback<Void>() {
            @Override
            public void onSuccess(Void value) {
                Log.d(TAG, "connect onSuccess");

                // Subscribe to a topic
                Topic[] topics = {new Topic(response.body().id, QoS.AT_LEAST_ONCE)};
                mqttConnection.subscribe(topics, new org.fusesource.mqtt.client.Callback<byte[]>() {
                    public void onSuccess(byte[] qoses) {
                        // The result of the subcribe request.
                        Log.d(TAG, "QoS: " + Arrays.toString(qoses));
                    }

                    public void onFailure(Throwable value) {
                        Log.d(TAG, "connect subscribe onFailure");
                        //mqttConnection.close(null); // subscribe failed.
                    }
                });
            }

            @Override
            public void onFailure(Throwable value) {
                Log.d(TAG, "connect onFailure");
            }
        });
    }

    private void mqttMsgHandler(byte[] payload) {
        Log.d(TAG, "Got a new msg via MQTT");

        // deserialize
        // {"type":"startNotification","msg":"testzz",id:1}
        // {"type":"stopNotification","msg":"testzz",id:1}
        String msgTxt = new String(payload, StandardCharsets.UTF_8);
        Gson gson = new Gson();
        Log.d(TAG, "Parsing JSON...");
        MqttMsg msgObj = null;
        try {
            msgObj = gson.fromJson(msgTxt, MqttMsg.class);
        } catch (Exception e) {
            Log.d(TAG, "Error when parsing JSON...");
            return;
        }

        Log.d(TAG, "JSON parsed correctly");
        // action based on msg
        if (msgObj.getType().equals("startNotification")) {
            Log.d(TAG, "Creating new notification...");
            startNotification(msgObj);
        } else if (msgObj.getType().equals("stopNotification")) {
            Log.d(TAG, "Removing notification...");
            stopNotification(msgObj);
        } else {
            Log.d(TAG, "Task in message not recognized...");
        }
    }

    private void stopNotification(MqttMsg msgObj) {
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
        notificationManager.cancel(msgObj.getSessionId());
    }

    private void startNotification(MqttMsg msgObj) {
        // prepare handler for clicking notification
        Intent taskView = new Intent(this, Task.class);
        taskView.putExtra("TASK_ID", msgObj.getId());
        taskView.putExtra("TASK_NAME", "...");
        // use System.currentTimeMillis() to have a unique ID for the pending intent
        PendingIntent pIntent = PendingIntent.getActivity(this, (int) System.currentTimeMillis(), taskView, 0);

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
                .setContentText("Counting...")
                .setUsesChronometer(true)
                .setWhen(currentTimeMillis())
                .setOngoing(true)
                .setContentIntent(pIntent)
                //.addAction(R.mipmap.ic_launcher, "Stop", stopActionIntent)
                .setColor(Color.GREEN)
                .setVibrate(new long[]{150})
                .setSound(Settings.System.DEFAULT_NOTIFICATION_URI)
                .setPriority(NotificationCompat.PRIORITY_DEFAULT);
        NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
        notificationManager.notify(msgObj.getSessionId(), builder.build());
    }

    private void foregroundNotification() {
        // Show notification to prevent killing in background
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel = new NotificationChannel("tasktab-service-running",
                    "Placeholder, disable me",
                    NotificationManager.IMPORTANCE_DEFAULT);
            NotificationManagerCompat notificationManager = NotificationManagerCompat.from(this);
            notificationManager.createNotificationChannel(channel);
        }
        final Intent notificationIntent = new Intent(this, MainActivity.class);
        final PendingIntent pi = PendingIntent.getActivity(this, 0, notificationIntent, 0);
        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, "tasktab-service-running")
                .setSmallIcon(R.drawable.ic_assignment_black_24dp)
                .setContentTitle("")
                //.setContentText("content text");
                .setOngoing(true)
                .setAutoCancel(false)
                .setContentIntent(pi);
        // for Android 7.1.1 let's hide notification icon because we can!
        if (Build.VERSION.SDK_INT == Build.VERSION_CODES.N_MR1) {
            builder.setPriority(NotificationCompat.PRIORITY_MIN);
        }
        final Notification notification = builder.build();
        startForeground(-1, notification);
    }
}
