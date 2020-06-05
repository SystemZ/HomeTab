package pl.systemz.tasktab;

import android.app.AlarmManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Intent;
import android.os.IBinder;
import android.util.Log;
import android.widget.Toast;

import me.pushy.sdk.Pushy;
import me.pushy.sdk.util.exceptions.PushyException;
import pl.systemz.tasktab.api.Client;
import pl.systemz.tasktab.receiver.Antenna;
import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;


public class StalkService extends Service {
    private static final String TAG = "StalkService";
    private Antenna antenna;

    private AlarmManager alarmMgr;
    private PendingIntent alarmIntent;

    IBinder mBinder;      // interface for clients that bind
    boolean mAllowRebind; // indicates whether onRebind should be used

    @Override
    public void onCreate() {
        // The service is being created
        Log.d(TAG, "onCreate()");
        Toast.makeText(this, "TaskStalk service started", Toast.LENGTH_LONG).show();

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
        //foregroundNotification();

        Thread thread = new Thread(() -> {
            try {
                // get device token from pushy.me for push notification
                Log.d(TAG, "trying to register push notifications");
                String pushToken = Pushy.register(getApplicationContext());
                Log.d(TAG, "Pushy device token: " + pushToken);

                // report token to backend
                Client client = Client.getInstance(getApplicationContext());
                Call<Void> call = client.getTtClient().deviceRegister(new Client.PushRegisterRequest(pushToken));
                call.enqueue(new Callback<Void>() {
                    @Override
                    public void onResponse(Call<Void> call, Response<Void> response) {
                        if (!response.isSuccessful()) {
                            Log.e(TAG, "push registration failed");
                            //justStop();
                            return;
                        }
                        Log.d(TAG, "push registartion ok");
                        //mqConnect(response);
                    }

                    @Override
                    public void onFailure(Call<Void> call, Throwable t) {
                        Log.e(TAG, "push registration onFailure()");
                        Log.e(TAG, t.getMessage());
                        //justStop();
                    }
                });
            } catch (PushyException e) {
                Log.e(TAG, e.toString());
            }
        });
        thread.start();

        // FIXME we shouldn't need this anymore
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
        Toast.makeText(this, "StalkTab Service stopped", Toast.LENGTH_LONG).show();
        Log.d(TAG, "onDestroy");
    }
}
