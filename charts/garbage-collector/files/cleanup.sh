  #!/bin/bash
  
  charts_to_delete=$(kubectl get pods -l "kube-tunnel=true" -o go-template='{{ range  $item := .items }}{{range .status.conditions }}{{ if (or (and (eq .type "PodScheduled") (eq .status "False")) (and (eq .type "Ready") (eq .status "False"))) }}{{index $item.metadata.labels "meta.helm.sh/release-name"}}{{"\n"}} {{ end }}{{ end }}{{ end }}' | grep -v "no value")
  
  if [ -n "$charts_to_delete" ]
  then
    echo "${charts_to_delete}" | xargs helm delete
  else
    echo "No charts to delete"
  fi

