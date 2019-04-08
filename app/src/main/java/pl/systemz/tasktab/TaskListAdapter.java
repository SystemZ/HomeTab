package pl.systemz.tasktab;

import android.content.Intent;
import android.text.TextUtils;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.View.OnClickListener;
import android.view.ViewGroup;
import android.widget.Filter;
import android.widget.Filterable;
import android.widget.TextView;

import java.util.ArrayList;
import java.util.List;

import androidx.recyclerview.widget.RecyclerView;
import pl.systemz.tasktab.api.Client;

public class TaskListAdapter extends RecyclerView.Adapter<TaskListAdapter.ViewHolder> implements Filterable {
    private List<Client.Timer> values;
    private List<Client.Timer> exampleList;
    private List<Client.Timer> exampleListFull;

    // Provide a reference to the views for each data item
    // Complex data items may need more than one view per item, and
    // you provide access to all the views for a data item in a view holder
    public class ViewHolder extends RecyclerView.ViewHolder {
        // each data item is just a string in this case
        public TextView txtHeader;
        public TextView txtFooter;
        public View layout;

        public ViewHolder(View v) {
            super(v);
            layout = v;
            txtHeader = v.findViewById(R.id.firstLine);
            txtFooter = v.findViewById(R.id.secondLine);
        }
    }

    public void add(int position, Client.Timer item) {
        values.add(position, item);
        notifyItemInserted(position);
    }

    public void remove(int position) {
        values.remove(position);
        notifyItemRemoved(position);
    }

    // Provide a suitable constructor (depends on the kind of dataset)
    public TaskListAdapter(List<Client.Timer> myDataset) {
        values = myDataset;
        this.exampleList = myDataset;
        exampleListFull = new ArrayList<>(exampleList);
    }

    // Create new views (invoked by the layout manager)
    @Override
    public TaskListAdapter.ViewHolder onCreateViewHolder(ViewGroup parent,
                                                         int viewType) {
        // create a new view
        LayoutInflater inflater = LayoutInflater.from(
                parent.getContext());
        View v = inflater.inflate(R.layout.row_layout, parent, false);
        // set the view's size, margins, paddings and layout parameters
        ViewHolder vh = new ViewHolder(v);
        return vh;
    }

    // Replace the contents of a view (invoked by the layout manager)
    @Override
    public void onBindViewHolder(ViewHolder holder, final int position) {
        // - get element from your dataset at this position
        // - replace the contents of the view with that element
        final Client.Timer task = values.get(position);
        // set title of this list element
        holder.txtHeader.setText(task.name);
        // set what to do on clicking
        holder.txtHeader.setOnClickListener(new OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent(v.getContext(), Task.class);
                intent.putExtra("TASK_ID", task.id);
                intent.putExtra("TASK_NAME", task.name);
                v.getContext().startActivity(intent);
            }
        });
        // set footer
        String footer = TextUtils.join(", ", task.tags);
        footer = footer + " (" + (task.seconds / 3600) + "h)";
        holder.txtFooter.setText(footer);
    }

    // Return the size of your dataset (invoked by the layout manager)
    @Override
    public int getItemCount() {
        return values.size();
    }

    @Override
    public Filter getFilter() {
        return exampleFilter;
    }

    private Filter exampleFilter = new Filter() {
        @Override
        protected FilterResults performFiltering(CharSequence constraint) {
            List<Client.Timer> filteredList = new ArrayList<>();
            if (constraint == null || constraint.length() == 0) {
                Log.v("", "zerolen");
                filteredList.addAll(exampleListFull);
            } else {
                String filterPattern = constraint.toString().toLowerCase().trim();
                Log.v("", "pattern: " + filterPattern);
                for (Client.Timer item : exampleListFull) {
//                    if (item.getText1().toLowerCase().startsWith(filterPattern)) {
                    if (item.name.toLowerCase().contains(filterPattern)) {
                        filteredList.add(item);
                    }
                }
            }

            FilterResults results = new FilterResults();
            results.values = filteredList;

            return results;
        }

        @Override
        protected void publishResults(CharSequence constraint, FilterResults results) {
            exampleList.clear();
            exampleList.addAll((List) results.values);
            notifyDataSetChanged();
        }
    };

}