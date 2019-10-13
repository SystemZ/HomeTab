package pl.systemz.tasktab.receiver;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.util.Log;

import pl.systemz.tasktab.StalkService;


public class BootReceiver extends BroadcastReceiver {
    private static final String TAG = "Stalk";

    @Override
    public void onReceive(Context context, Intent intent) {
        Log.d(TAG, "BootReceiver onReceive, starting service");
        Intent serviceIntent = new Intent(context, StalkService.class);
        context.startService(serviceIntent);
    }
}
