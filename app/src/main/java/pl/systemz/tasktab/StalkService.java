package pl.systemz.tasktab;

import android.app.AlarmManager;
import android.app.PendingIntent;
import android.app.Service;
import android.content.Intent;
import android.content.IntentFilter;
import android.os.IBinder;
import android.util.Log;
import android.widget.Toast;

public class StalkService extends Service {
    private static final String TAG = "Stalk";
    private Antenna antenna;

    private AlarmManager alarmMgr;
    private PendingIntent alarmIntent;

    int mStartMode;       // indicates how to behave if the service is killed
    IBinder mBinder;      // interface for clients that bind
    boolean mAllowRebind; // indicates whether onRebind should be used

    @Override
    public void onCreate() {
        // The service is being created

        //Intent intents = new Intent(getBaseContext(),MainActivity.class);
        //intents.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        //startActivity(intents);

        Toast.makeText(this, "Stalk service started", Toast.LENGTH_LONG).show();
        Log.d(TAG, "onStart, register events");

        // Hopefully your alarm will have a lower frequency than this!
//        AlarmManager alarmMgr;
//        alarmMgr = (AlarmManager)this.getSystemService(Context.ALARM_SERVICE);
//        pi = PendingIntent.getBroadcast(this, 11, new Intent("pl.systemz.stalk.GPSSETTINGS"), 0);
//        am.setInexactRepeating(AlarmManager.RTC_WAKEUP, 0, TWENTYFIVE_SECONDS, pi);

//        alarmMgr.setInexactRepeating(AlarmManager.ELAPSED_REALTIME_WAKEUP,
//                SystemClock.elapsedRealtime() + 60*1000,
//                60*1000, alarmIntent);

//                SystemClock.elapsedRealtime() + AlarmManager.INTERVAL_HALF_HOUR,
//                AlarmManager.INTERVAL_HALF_HOUR, alarmIntent);


        IntentFilter intentFilter = new IntentFilter();
        intentFilter.addAction(Intent.ACTION_SCREEN_ON);
        intentFilter.addAction(Intent.ACTION_SCREEN_OFF);
        intentFilter.addAction(Intent.ACTION_BATTERY_CHANGED);
        intentFilter.addAction(Intent.ACTION_BOOT_COMPLETED);

        this.antenna = new Antenna();
        registerReceiver(antenna, intentFilter);

    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        // The service is starting, due to a call to startService()
        Log.d(TAG, "onStartCommand");

        return mStartMode;
    }

    @Override
    public IBinder onBind(Intent intent) {
        // A client is binding to the service with bindService()
        return mBinder;
    }

    @Override
    public boolean onUnbind(Intent intent) {
        // All clients have unbound with unbindService()
        return mAllowRebind;
    }

    @Override
    public void onRebind(Intent intent) {
        // A client is binding to the service with bindService(),
        // after onUnbind() has already been called
    }

    @Override
    public void onDestroy() {
        // The service is no longer used and is being destroyed
        Toast.makeText(this, "Stalk Service stopped", Toast.LENGTH_LONG).show();
        Log.d(TAG, "onDestroy");
    }

}
