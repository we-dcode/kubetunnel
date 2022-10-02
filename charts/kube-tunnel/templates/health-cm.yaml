apiVersion: v1
kind: Secret
metadata:
  name: {{ include "kubetunnel.fullname" . }}-{{ .Values.env_service_name }}-health
  labels:
  {{- include "kubetunnel.labels" . | nindent 4 }}
stringData:
  health.sh: |
    {{.Values.healthInline | nindent 4}}
