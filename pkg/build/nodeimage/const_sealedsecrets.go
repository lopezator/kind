package nodeimage

var defaultSealedSecretsImages = []string{"quay.io/bitnami/sealed-secrets-controller:v0.9.5"}

const defaultSealedSecretsManifest = `
# sealed-secrets manifest
# would you kindly template this file
---
apiVersion: v1
items:
- apiVersion: v1
  kind: Namespace
  metadata:
    name: sealed-secrets
- apiVersion: apiextensions.k8s.io/v1beta1
  kind: CustomResourceDefinition
  metadata:
    name: sealedsecrets.bitnami.com
  spec:
    group: bitnami.com
    names:
      kind: SealedSecret
      listKind: SealedSecretList
      plural: sealedsecrets
      singular: sealedsecret
    scope: Namespaced
    version: v1alpha1
- apiVersion: v1
  kind: Service
  metadata:
    annotations: {}
    labels:
      name: sealed-secrets-controller
    name: sealed-secrets-controller
    namespace: sealed-secrets
  spec:
    ports:
    - port: 8080
      targetPort: 8080
    selector:
      name: sealed-secrets-controller
    type: ClusterIP
- apiVersion: rbac.authorization.k8s.io/v1beta1
  kind: RoleBinding
  metadata:
    annotations: {}
    labels:
      name: sealed-secrets-service-proxier
    name: sealed-secrets-service-proxier
    namespace: sealed-secrets
  roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: Role
    name: sealed-secrets-service-proxier
  subjects:
  - apiGroup: rbac.authorization.k8s.io
    kind: Group
    name: system:authenticated
- apiVersion: rbac.authorization.k8s.io/v1beta1
  kind: Role
  metadata:
    annotations: {}
    labels:
      name: sealed-secrets-service-proxier
    name: sealed-secrets-service-proxier
    namespace: sealed-secrets
  rules:
  - apiGroups:
    - ''
    resourceNames:
    - 'http:sealed-secrets-controller:'
    - sealed-secrets-controller
    resources:
    - services/proxy
    verbs:
    - create
    - get
- apiVersion: rbac.authorization.k8s.io/v1beta1
  kind: RoleBinding
  metadata:
    annotations: {}
    labels:
      name: sealed-secrets-controller
    name: sealed-secrets-controller
    namespace: sealed-secrets
  roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: Role
    name: sealed-secrets-key-admin
  subjects:
  - kind: ServiceAccount
    name: sealed-secrets-controller
    namespace: sealed-secrets
- apiVersion: rbac.authorization.k8s.io/v1beta1
  kind: Role
  metadata:
    annotations: {}
    labels:
      name: sealed-secrets-key-admin
    name: sealed-secrets-key-admin
    namespace: sealed-secrets
  rules:
  - apiGroups:
    - ''
    resources:
    - secrets
    verbs:
    - create
    - list
- apiVersion: v1
  kind: ServiceAccount
  metadata:
    annotations: {}
    labels:
      name: sealed-secrets-controller
    name: sealed-secrets-controller
    namespace: sealed-secrets
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    annotations: {}
    labels:
      name: sealed-secrets-controller
    name: sealed-secrets-controller
    namespace: sealed-secrets
  spec:
    minReadySeconds: 30
    replicas: 1
    revisionHistoryLimit: 10
    selector:
      matchLabels:
        name: sealed-secrets-controller
    strategy:
      rollingUpdate:
        maxSurge: 25%
        maxUnavailable: 25%
      type: RollingUpdate
    template:
      metadata:
        annotations: {}
        labels:
          name: sealed-secrets-controller
      spec:
        containers:
        - args: []
          command:
          - controller
          - "--key-renew-period=0"
          env: []
          image: quay.io/bitnami/sealed-secrets-controller:v0.9.5
          imagePullPolicy: Always
          livenessProbe:
            httpGet:
              path: "/healthz"
              port: http
          name: sealed-secrets-controller
          ports:
          - containerPort: 8080
            name: http
          readinessProbe:
            httpGet:
              path: "/healthz"
              port: http
          securityContext:
            readOnlyRootFilesystem: true
            runAsNonRoot: true
            runAsUser: 1001
          stdin: false
          tty: false
          volumeMounts:
          - mountPath: "/tmp"
            name: tmp
        imagePullSecrets: []
        initContainers: []
        serviceAccountName: sealed-secrets-controller
        terminationGracePeriodSeconds: 30
        volumes:
        - emptyDir: {}
          name: tmp
- apiVersion: rbac.authorization.k8s.io/v1beta1
  kind: ClusterRoleBinding
  metadata:
    annotations: {}
    labels:
      name: sealed-secrets-controller
    name: sealed-secrets-controller
  roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: secrets-unsealer
  subjects:
  - kind: ServiceAccount
    name: sealed-secrets-controller
    namespace: sealed-secrets
- apiVersion: rbac.authorization.k8s.io/v1beta1
  kind: ClusterRole
  metadata:
    annotations: {}
    labels:
      name: secrets-unsealer
    name: secrets-unsealer
  rules:
  - apiGroups:
    - bitnami.com
    resources:
    - sealedsecrets
    verbs:
    - get
    - list
    - watch
    - update
  - apiGroups:
    - ''
    resources:
    - secrets
    verbs:
    - get
    - create
    - update
    - delete
  - apiGroups:
    - ''
    resources:
    - events
    verbs:
    - create
    - patch
kind: List
`
