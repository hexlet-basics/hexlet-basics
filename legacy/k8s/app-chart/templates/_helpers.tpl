{{/*
Общие фрагменты чарта codebasics. Вынесены из шаблонов, чтобы убрать копипасту.
Важно: рендер должен оставаться идентичным прежнему — селекторы и pod-template
labels НЕ трогаем (иммутабельны у Deployment/DaemonSet).
*/}}

{{/* Общие метки верхнего уровня (metadata.labels) ресурсов чарта. */}}
{{- define "codebasics.labels" -}}
app.kubernetes.io/name: "{{ .Chart.Name }}"
app.kubernetes.io/instance: "{{ .Release.Name }}"
app.kubernetes.io/version: "{{ .Chart.Version }}"
app.kubernetes.io/managed-by: "{{ .Release.Service }}"
{{- end -}}

{{/* nodeSelector + imagePullSecrets — общий блок планирования подов. */}}
{{- define "codebasics.podScheduling" -}}
nodeSelector:
  group: "{{ .Release.Namespace }}"
imagePullSecrets:
  - name: docker-config
{{- end -}}

{{/* nodeAffinity: планирование на ноды роли service. */}}
{{- define "codebasics.nodeAffinityService" -}}
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
        - matchExpressions:
            - key: role
              operator: In
              values:
                - service
{{- end -}}

{{/* Тома для контейнеров, запускающих код учеников (hostPath + docker secret). */}}
{{- define "codebasics.appVolumes" -}}
volumes:
  - name: users-code-path
    hostPath:
      path: /tmp/hexlet-basics
  - name: docker-sock
    hostPath:
      path: /var/run/docker.sock
  - name: docker-secret
    secret:
      secretName: docker-config
      items:
        - key: .dockerconfigjson
          path: config.json
{{- end -}}
