package pl.systemz.hometab.receiver;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.util.Log;


public class Antenna extends BroadcastReceiver {
    private static final String TAG = "Stalk";

    @Override
    public void onReceive(Context context, Intent intent) {
        Log.e(TAG, "onReceive:");
        Log.d(TAG, intent.getAction());
    }
}
